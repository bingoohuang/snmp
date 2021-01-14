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
	Targets   arrayFlags
	Oids      arrayFlags

	TrapAddr string
	Mode     string
	Logger   *log.Logger
}

type arrayFlags []string

func (i *arrayFlags) String() string { return "my string representation" }

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (o *Options) ParseFlags() {
	var x, y, z arrayFlags

	flag.StringVar(&o.Mode, "mode", "get/walk", "")
	flag.StringVar(&o.Community, "c", "public", "")
	flag.Var(&o.Targets, "t", "")
	flag.Var(&x, "x", "")
	flag.Var(&y, "y", "")
	flag.Var(&z, "z", "")
	flag.Var(&o.Oids, "oid", "")
	flag.StringVar(&o.TrapAddr, "trap", "", "")
	verbose := flag.Bool("V", false, "")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `Usage of snmp: snmp [options] Oids...
  -mode  get/walk/trapsend (default is get/walk)
  -c     string Default SNMP community (default "public")
  -t     one or more SNMP targets (eg. -t 192.168.1.1 -t myCommunity@192.168.1.2:1234)
  -x/y/z one or more x/y/z vars (eg. -x 1-3)
  -oids  one or more Oids
  -trap  trap server listening address(eg. :9162)
  -V     Verbose logging of packets
`)
	}

	flag.Parse()

	o.Oids = append(o.Oids, flag.Args()...)
	o.Oids = interpolate(o.Oids, expandNums(x), "x")
	o.Oids = interpolate(o.Oids, expandNums(y), "y")
	o.Oids = interpolate(o.Oids, expandNums(z), "z")

	if *verbose {
		o.Logger = log.New(log.Writer(), log.Prefix(), log.Flags())
	}
}

func interpolate(args []string, xs []string, xName string) []string {
	vs := make([]string, 0)

	for _, arg := range args {
		if !strings.Contains(arg, "."+xName) {
			vs = append(vs, arg)
			continue
		}

		for _, x := range xs {
			y := strings.ReplaceAll(arg, "."+xName, "."+x)
			vs = append(vs, y)
		}
	}

	return vs
}

type ExpandNums struct {
	nums []string
}

func expandNums(x []string) []string {
	ox := ExpandNums{nums: make([]string, 0)}
	for _, xi := range x {
		ox.expand(xi)
	}

	return ox.nums
}

func (ox *ExpandNums) expand(xi string) {
	xii := strings.Split(xi, ",")
	for _, xij := range xii {
		if !strings.Contains(xij, "-") {
			ox.append(xij)
			continue
		}

		xirange := strings.Split(xij, "-")
		f, err := strconv.Atoi(xirange[0])
		if err != nil {
			log.Printf("W! error x values %v", err)
			continue
		}

		ox.append(xirange[0])

		xirange = xirange[1:]
		if len(xirange) == 0 {
			continue
		}

		ox.expandRange(f, xirange)
	}
}

func (ox *ExpandNums) append(f string) {
	for _, i := range ox.nums {
		if i == f {
			return
		}
	}

	ox.nums = append(ox.nums, f)
}

func (ox *ExpandNums) expandRange(f int, to []string) {
	t, err := strconv.Atoi(to[0])
	if err != nil {
		log.Printf("W! error x values %v", err)
		return
	}

	to = to[1:]
	step := 1
	if len(to) > 0 {
		v, err := strconv.Atoi(to[0])
		if err != nil {
			log.Printf("W! error x values %v", err)
			return
		}
		step = v
	}

	for j := f + step; j <= t; j += step {
		ox.append(fmt.Sprintf("%d", j))
	}
}

func main() {
	options := Options{}
	options.ParseFlags()

	for _, t := range options.Targets {
		options.do(t)
	}

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

func (t *Target) snmpWalk() {
	if !strings.Contains(t.Mode, "walk") {
		return
	}

	for _, oid := range t.Oids {
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

	result, err := t.Get(t.Oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		log.Printf("W! snmpget %v error: %v", t.Oids, err)
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
	fmt.Printf("[%s][%s][%d] %s = %v: ", typ, target, i, pdu.Name, pdu.Type)

	switch pdu.Type {
	case g.OctetString:
		fmt.Printf("%s\n", pdu.Value.([]byte))
	case g.ObjectIdentifier:
		fmt.Printf("%s\n", pdu.Value.(string))
	default:
		fmt.Printf("%v\n", pdu.Value)
	}
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
	tl.Params.Logger = o.Logger

	if err := tl.Listen(o.TrapAddr); err != nil {
		log.Printf("E! error in listen: %s", err)
	}
}

func trapHandler(p *g.SnmpPacket, addr *net.UDPAddr) {
	log.Printf("got trapdata from %s", addr.IP)
	for i, v := range p.Variables {
		printPdu("trap", addr.String(), i, v)
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
