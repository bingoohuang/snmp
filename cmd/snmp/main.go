package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bingoohuang/snmp/pkg/snmpp"
	g "github.com/gosnmp/gosnmp"
)

func main() {
	options := Options{}
	options.ParseFlags()

	for _, t := range options.Agents {
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
	t := o.createAgent(target)
	if err := t.Connect(); err != nil {
		log.Printf("E! Connect() err: %v", err)
		os.Exit(1)
	}

	defer t.Close()

	t.trapSend()
	t.snmpGet()
	t.snmpWalk()
}

func (o *Options) createAgent(target string) Target {
	gs := o.CreateGoSNMP(o.Logger != nil)
	if o.Logger != nil {
		gs.Logger = *o.Logger
	}

	target = o.parseCommunity(target, gs)
	if err := snmpp.SetTarget(gs, target); err != nil {
		log.Fatalf("set target failed: %v", err)
	}

	return Target{
		GoSNMP:  gs,
		target:  refineTargetForOutput(gs),
		Options: o,
	}
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

func refineTargetForOutput(gs *g.GoSNMP) string {
	target := ""

	if gs.Community != snmpp.DefaultCommunity {
		target = gs.Community + "@"
	}

	target += gs.Target
	if gs.Port != snmpp.DefaultSnmpPort {
		target += fmt.Sprintf(":%d", gs.Port)
	}
	return target
}
