package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bingoohuang/snmp/pkg/smi"
	"github.com/bingoohuang/snmp/pkg/snmpp"
	"github.com/bingoohuang/snmp/pkg/util"
)

type Options struct {
	Community string
	Targets   arrayFlags
	Oids      arrayFlags

	TrapAddr string
	Mode     string
	Logger   *log.Logger
	mib      *smi.MIB
}

type arrayFlags []string

func (i *arrayFlags) String() string { return "my string representation" }

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (o *Options) ParseFlags() {
	var x, y, z arrayFlags

	flag.StringVar(&o.Mode, "m", "get/walk", "")
	flag.StringVar(&o.Community, "c", "public", "")
	flag.Var(&o.Targets, "t", "")
	flag.Var(&x, "x", "")
	flag.Var(&y, "y", "")
	flag.Var(&z, "z", "")
	flag.Var(&o.Oids, "oid", "")
	flag.StringVar(&o.TrapAddr, "s", "", "")
	verbose := flag.Bool("V", false, "")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `Usage of snmp: snmp [options] Oids...
  -m     get/walk/trapsend/translate (default is get/walk)
  -c     string Default SNMP community (default "public")
  -t     one or more SNMP targets (eg. -t 192.168.1.1 -t myCommunity@192.168.1.2:1234)
  -x/y/z one or more x/y/z vars (eg. -x 1-3)
  -oid   one or more Oids
  -s     trap server listening address(eg. :9162)
  -V     Verbose logging of packets
`)
	}

	flag.Parse()

	o.mib = snmpp.LoadMibs()
	o.Oids = append(o.Oids, flag.Args()...)
	isTranslate := o.Mode == "translate"
	o.Oids = interpolate(isTranslate, o.mib, o.Oids, util.ExpandNums(x), "x")
	o.Oids = interpolate(isTranslate, o.mib, o.Oids, util.ExpandNums(y), "y")
	o.Oids = interpolate(isTranslate, o.mib, o.Oids, util.ExpandNums(z), "z")

	if *verbose {
		o.Logger = log.New(log.Writer(), log.Prefix(), log.Flags())
		log.Printf("Oids:%v", o.Oids)
	}
}

func interpolate(isTranslate bool, mib *smi.MIB, args []string, xs []string, xName string) []string {
	vs := make([]string, 0)

	for _, arg := range args {
		if snmpp.IsSymbolName(arg) && !isTranslate {
			oid, err := mib.OID(arg)
			if err != nil {
				log.Printf("unkown symbol %s", arg)
				continue
			}
			vs = append(vs, oid.String())
		}

		if snmpp.IsSymbolName(arg) || !strings.Contains(arg, "."+xName) {
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
