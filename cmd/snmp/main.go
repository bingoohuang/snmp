package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	g "github.com/gosnmp/gosnmp"
)

type Options struct {
	Community string
	Verbose   bool
	Targets   arrayFlags
	TrapAddr  string
	Mode      string
}

type arrayFlags []string

func (i *arrayFlags) String() string { return "my string representation" }

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (o *Options) InitFlags() {
	flag.StringVar(&o.Community, "c", "public", "")
	flag.StringVar(&o.Mode, "mode", "get/walk", "")
	flag.BoolVar(&o.Verbose, "V", false, "")
	flag.Var(&o.Targets, "t", "")
	flag.StringVar(&o.TrapAddr, "trap", "", "")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `Usage of snmp:
  -c    string Default SNMP community (default "public")
  -mode get/walk/trapsend (default is get/walk)
  -t    value Default SNMP community
  -trap trap server listening address, eg: :9162
  -V    Verbose logging of packets
`)
	}
}

func main() {
	options := Options{}
	options.InitFlags()

	flag.Parse()

	for _, t := range options.Targets {
		options.do(t, flag.Args())
	}

	options.trap()
}

type Target struct {
	*g.GoSNMP
	*Options
	target string
	oids   []string
}

func (t *Target) Close() {
	_ = t.Conn.Close()
}

func (o *Options) do(target string, oids []string) {
	t := o.createTarget(target, oids)
	if err := t.Connect(); err != nil {
		log.Printf("E! Connect() err: %v", err)
		os.Exit(1)
	}

	defer t.Close()

	t.trapSend()
	t.snmpGet()
	t.snmpWalk()
}

func (t *Target) snmpWalk() {
	if !strings.Contains(t.Mode, "walk") {
		return
	}

	for _, oid := range t.oids {
		i := 0
		if err := t.BulkWalk(oid, func(pdu g.SnmpPDU) error {
			printPdu("walk", t.target, i, pdu)
			i++
			return nil
		}); err != nil {
			log.Printf("W! snmpwalk error: %v", err)
		}
	}
}

func (t *Target) snmpGet() {
	if !strings.Contains(t.Mode, "get") {
		return
	}

	result, err := t.Get(t.oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		log.Printf("W! snmpget error: %v", err)
		return
	}

	for i, pdu := range result.Variables {
		printPdu("get", t.target, i, pdu)
	}
}

func (t *Target) trapSend() {
	if !strings.Contains(t.Mode, "trapsend") {
		return
	}

	trap := g.SnmpTrap{
		Variables: []g.SnmpPDU{{
			Name:  "1.3.6.1.2.1.1.6",
			Type:  g.ObjectIdentifier,
			Value: "1.3.6.1.2.1.1.6.10",
		}, {
			Name:  "1.3.6.1.2.1.1.7",
			Type:  g.OctetString,
			Value: "Testing TCP trap...",
		}, {
			Name:  "1.3.6.1.2.1.1.8",
			Type:  g.Integer,
			Value: 123,
		}},
	}

	if _, err := t.SendTrap(trap); err != nil {
		log.Printf("E! SendTrap() err: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func printPdu(typ, target string, i int, pdu g.SnmpPDU) {
	fmt.Printf("[%s][%s][%d] %s = ", typ, target, i, pdu.Name)

	switch pdu.Type {
	case g.OctetString:
		fmt.Printf("%v: %s\n", pdu.Type, pdu.Value.([]byte))
	case g.ObjectIdentifier:
		fmt.Printf("%v: %s\n", pdu.Type, pdu.Value.(string))
	default:
		fmt.Printf("%v: %v\n", pdu.Type, pdu.Value)
	}
}

const (
	DefaultSnmpPort = 161
)

func (o *Options) createTarget(target string, oids []string) Target {
	gs := &g.GoSNMP{
		Port:               DefaultSnmpPort,
		Transport:          "udp",
		Community:          "public",
		Version:            g.Version2c,
		Timeout:            time.Duration(10) * time.Second,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            g.MaxOids,
	}

	o.setupVerbose(gs)
	target = o.parseCommunity(target, gs)
	o.parseTargetPort(target, gs)
	refinedTarget := refineTargetForOutput(gs)

	return Target{
		GoSNMP:  gs,
		target:  refinedTarget,
		oids:    oids,
		Options: o,
	}
}

func (o *Options) setupVerbose(gs *g.GoSNMP) {
	if !o.Verbose {
		return
	}

	gs.Logger = log.New(log.Writer(), log.Prefix(), log.Flags())

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

func (o *Options) parseTargetPort(target string, gs *g.GoSNMP) {
	if p := strings.LastIndex(target, ":"); p < 0 {
		gs.Target = target
	} else {
		gs.Target = target[:p]
		port, err := strconv.ParseUint(target[p+1:], 10, 16)
		if err != nil {
			log.Fatalf("parse port error %s: %v", target, err)
		}
		gs.Port = uint16(port)
	}
}

func (o *Options) trap() {
	if o.TrapAddr == "" {
		return
	}

	tl := g.NewTrapListener()
	tl.OnNewTrap = trapHandler
	tl.Params = g.Default

	if o.Verbose {
		tl.Params.Logger = log.New(log.Writer(), log.Prefix(), log.Flags())
	}

	if err := tl.Listen(o.TrapAddr); err != nil {
		log.Printf("E! error in listen: %s", err)
	}
}

func trapHandler(packet *g.SnmpPacket, addr *net.UDPAddr) {
	log.Printf("got trapdata from %s", addr.IP)
	for i, v := range packet.Variables {
		printPdu("trap", addr.String(), i, v)
	}
}

func refineTargetForOutput(gs *g.GoSNMP) string {
	target := ""

	if gs.Community != "public" {
		target = gs.Community + "@"
	}

	target += gs.Target
	if gs.Port != DefaultSnmpPort {
		target += fmt.Sprintf(":%d", gs.Port)
	}
	return target
}
