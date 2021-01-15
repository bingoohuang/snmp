package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	g "github.com/gosnmp/gosnmp"
)

func main() {
	options := Options{}
	options.ParseFlags()

	for _, t := range options.Targets {
		options.do(t)
	}

	options.translate()
	options.trap()
}

type Target struct {
	*g.GoSNMP
	*Options
	target string
}

func (t *Target) Close() {
	_ = t.Conn.Close()
}

func (o *Options) do(target string) {
	t := o.createTarget(target)
	if err := t.Connect(); err != nil {
		log.Printf("E! Connect() err: %v", err)
		os.Exit(1)
	}

	defer t.Close()

	t.trapSend()
	t.snmpGet()
	t.snmpWalk()
}

const (
	DefaultSnmpPort  = 161
	DefaultCommunity = "public"
)

func (o *Options) createTarget(target string) Target {
	gs := &g.GoSNMP{
		Port:               DefaultSnmpPort,
		Transport:          "udp",
		Community:          DefaultCommunity,
		Version:            g.Version2c,
		Timeout:            time.Duration(3) * time.Second,
		Retries:            0,
		ExponentialTimeout: false,
		MaxOids:            g.MaxOids,
	}

	o.setupVerbose(gs)
	target = o.parseCommunity(target, gs)
	parseTargetPort(target, gs)
	refinedTarget := refineTargetForOutput(gs)

	return Target{
		GoSNMP:  gs,
		target:  refinedTarget,
		Options: o,
	}
}

func (o *Options) setupVerbose(gs *g.GoSNMP) {
	if o.Logger == nil {
		return
	}

	gs.Logger = o.Logger

	// Function handles for collecting metrics on query latencies.
	var sent time.Time
	gs.OnSent = func(*g.GoSNMP) { sent = time.Now() }
	gs.OnRecv = func(*g.GoSNMP) { log.Printf("Query latency: %s", time.Since(sent)) }
}

func (o *Options) parseCommunity(target string, gs *g.GoSNMP) string {
	if p := strings.LastIndex(target, "@"); p < 0 {
		gs.Community = o.Community
	} else {
		gs.Community = target[:p]
		target = target[p+1:]
	}

	return target
}

func parseTargetPort(target string, gs *g.GoSNMP) {
	if p := strings.LastIndex(target, ":"); p < 0 {
		gs.Target = target
		return
	} else {
		gs.Target = target[:p]
		port, err := strconv.ParseUint(target[p+1:], 10, 16)
		if err != nil {
			log.Fatalf("parse port error %s: %v", target, err)
		}
		gs.Port = uint16(port)
	}
}

func refineTargetForOutput(gs *g.GoSNMP) string {
	target := ""

	if gs.Community != DefaultCommunity {
		target = gs.Community + "@"
	}

	target += gs.Target
	if gs.Port != DefaultSnmpPort {
		target += fmt.Sprintf(":%d", gs.Port)
	}
	return target
}
