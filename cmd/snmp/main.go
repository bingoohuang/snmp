package main

import (
	"fmt"
	"log"
	"os"

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

	pduNames := map[string]bool{}

	t.snmpGet(pduNames)
	t.snmpWalk(pduNames)
}

func (o *Options) createAgent(target string) Target {
	gs := o.CreateGoSNMP(target, o.Logger != nil)
	if o.Logger != nil {
		gs.Logger = *o.Logger
	}

	return Target{
		GoSNMP:  gs,
		target:  refineTargetForOutput(gs),
		Options: o,
	}
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
