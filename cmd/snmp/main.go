package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	g "github.com/gosnmp/gosnmp"
)

type Options struct {
	Community string
	Verbose   bool
	Targets   arrayFlags
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (o *Options) InitFlags() {
	flag.StringVar(&o.Community, "c", "public", "Default SNMP community")
	flag.BoolVar(&o.Verbose, "V", false, "Do verbose logging of packets")
	flag.Var(&o.Targets, "t", "Default SNMP community")
}

func main() {
	options := Options{}
	options.InitFlags()

	flag.Parse()

	for _, t := range options.Targets {
		options.do(t, flag.Args())
	}
}

func (o *Options) do(target string, oids []string) {
	target, gs, err := o.createSnmp(target)
	if err != nil {
		log.Printf("W! failed to create snmp error: %v", err)
		return
	}

	if err := gs.Connect(); err != nil {
		log.Printf("E! Connect() err: %v", err)
		return
	}

	defer gs.Conn.Close()

	if err := snmpGet(target, gs, oids); err != nil {
		log.Printf("W! snmpget error: %v", err)
	}

	for _, oid := range oids {
		i := 0
		if err := gs.BulkWalk(oid, func(pdu g.SnmpPDU) error {
			printPdu("walk", target, i, pdu)
			i++
			return nil
		}); err != nil {
			log.Printf("W! snmpget walk: %v", err)
		}
	}
}

func snmpGet(target string, gs *g.GoSNMP, oids []string) error {
	result, err := gs.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		return err
	}

	for i, pdu := range result.Variables {
		printPdu(" get", target, i, pdu)
	}

	return nil
}

func printPdu(typ, target string, i int, pdu g.SnmpPDU) {
	fmt.Printf("[%s][%s][%d] %s = ", typ, target, i, pdu.Name)

	switch pdu.Type {
	case g.OctetString:
		fmt.Printf("string: %s\n", pdu.Value.([]byte))
	default:
		// ... or often you're just interested in numeric values.
		// ToBigInt() will return the Value as a BigInt, for plugging
		// into your calculations.
		fmt.Printf("number: %d\n", g.ToBigInt(pdu.Value))
	}
}

const (
	DefaultSnmpPort = 161
)

func (o *Options) createSnmp(target string) (refinedTarget string, gs *g.GoSNMP, err error) {
	gs = &g.GoSNMP{
		Port:               DefaultSnmpPort,
		Transport:          "udp",
		Community:          "public",
		Version:            g.Version2c,
		Timeout:            time.Duration(10) * time.Second,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            g.MaxOids,
	}

	if o.Verbose {
		gs.Logger = log.New(log.Writer(), log.Prefix(), log.Flags())

		// Function handles for collecting metrics on query latencies.
		var sent time.Time
		gs.OnSent = func(*g.GoSNMP) { sent = time.Now() }
		gs.OnRecv = func(*g.GoSNMP) { log.Printf("Query latency: %s", time.Since(sent)) }
	}

	if p := strings.LastIndex(target, "@"); p < 0 {
		gs.Community = o.Community
	} else {
		gs.Community = target[:p]
		target = target[p+1:]
	}

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

	if gs.Community != "public" {
		refinedTarget = gs.Community + "@"
	}

	refinedTarget += gs.Target
	if gs.Port != DefaultSnmpPort {
		refinedTarget += fmt.Sprintf(":%d", gs.Port)
	}

	return
}
