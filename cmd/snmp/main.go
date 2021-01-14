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
	flag.Var(&o.Targets, "t", "Default SNMP community")
}

func main() {
	options := Options{}
	options.InitFlags()

	flag.Parse()

	for _, t := range options.Targets {
		do(t, options, flag.Args())
	}
}

func do(target string, options Options, oids []string) {
	target, gs, err := createSnmp(target, options)
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
		log.Printf("E! Get() err: %v", err)
	}
}

func snmpGet(target string, gs *g.GoSNMP, oids []string) error {
	result, err := gs.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		return err
	}

	for i, v := range result.Variables {
		fmt.Printf("[%s] [%d] %s = ", target, i, v.Name)

		switch v.Type {
		case g.OctetString:
			fmt.Printf("string: %s\n", v.Value.([]byte))
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			fmt.Printf("number: %d\n", g.ToBigInt(v.Value))
		}
	}

	return nil
}

const (
	DefaultSnmpPort = 161
)

func createSnmp(target string, options Options) (refinedTarget string, gs *g.GoSNMP, err error) {
	gs = &g.GoSNMP{
		Port:               DefaultSnmpPort,
		Transport:          "udp",
		Community:          "public",
		Version:            g.Version2c,
		Timeout:            time.Duration(2) * time.Second,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            g.MaxOids,
	}

	if p := strings.LastIndex(target, "@"); p < 0 {
		gs.Community = options.Community
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
