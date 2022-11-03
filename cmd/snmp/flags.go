package main

import (
	"flag"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/bingoohuang/snmp/pkg/smi"
	"github.com/bingoohuang/snmp/pkg/snmpp"
	"github.com/bingoohuang/snmp/pkg/util"
)

type ClientConfig struct {
	Timeout time.Duration `usage:"Timeout to wait for a response."`

	Retries              int    `usage:"Number of retries to attempt."`
	Version              int    `usage:"SNMP version; can be 1, 2, or 3."`
	UnconnectedUDPSocket bool   `usage:"Unconnected UDP socket When true, SNMP responses are accepted from any address not just the requested address. This can be useful when gathering from redundant/failover systems."`
	Community            string `usage:"SNMP community string, Parameters for Version 1 & 2"`

	MaxRepetitions int `usage:"The GETBULK max-repetitions parameter, Parameters for Version 2 & 3"`

	// Parameters for Version 3
	// SNMPv3 authentication and encryption options.

	ContextName  string `usage:"Context Name. for SNMPv3"`
	SecLevel     string `usage:"Security Level; one of noAuthNoPriv, authNoPriv, or authPriv. for SNMPv3"`
	SecName      string `usage:"Security Name. for SNMPv3"`
	AuthProtocol string `usage:"Authentication protocol; one of MD5, SHA, SHA224, SHA256, SHA384, SHA512. for SNMPv3"`
	AuthPassword string `usage:"Authentication password"`
	// Protocols "AES192", "AES192", "AES256", and "AES256C" require the underlying net-snmp tools
	// to be compiled with --enable-blumenthal-aes (http://www.net-snmp.org/docs/INSTALL.html)
	PrivProtocol string `usage:"Privacy protocol used for encrypted messages; one of DES, AES, AES192, AES192C, AES256, AES256C"`
	PrivPassword string `usage:"Privacy password used for encrypted messages"`

	EngineID    string
	EngineBoots int
	EngineTime  int
}

type Options struct {
	ClientConfig

	Agents  arrayFlags
	Oids    arrayFlags
	Verbose string

	TrapAddr string
	Operate  string
	Logger   *gosnmp.Logger
	mib      *smi.MIB
}

type arrayFlags []string

func (i *arrayFlags) String() string { return "my string representation" }

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	timeDurationType = reflect.TypeOf(time.Duration(0))
)

type ArrayFlags struct {
	Value string
	pp    *[]string
}

func (i *ArrayFlags) String() string { return i.Value }

func (i *ArrayFlags) Set(value string) error {
	*i.pp = append(*i.pp, value)
	return nil
}

func (o *Options) ParseFlags() {
	var x, y, z arrayFlags

	flag.StringVar(&o.Operate, "m", "get/walk", "")
	flag.Var(&o.Agents, "t", `Agent addresses to retrieve values from.`+"\n"+
		`   format:  agents = ["<scheme://><hostname>:<optional port>"]`+"\n"+
		`   scheme:  optional, either udp, udp4, udp6, tcp, tcp4, tcp6. default is udp`+"\n"+
		`   example: 127.0.0.1, udp://127.0.0.1:161, tcp://127.0.0.1:161, udp4://v4only-snmp-agent`+"\n")
	flag.Var(&x, "x", "")
	flag.Var(&y, "y", "")
	flag.Var(&z, "z", "")
	flag.Var(&o.Oids, "o", "")
	flag.StringVar(&o.TrapAddr, "s", "", "")
	flag.StringVar(&o.Verbose, "V", "", "debug,desc")

	ccValue := reflect.ValueOf(&o.ClientConfig).Elem()
	ccType := ccValue.Type()
	var usages []string
	for i := 0; i < ccType.NumField(); i++ {
		ft := ccType.Field(i)
		if !ft.IsExported() {
			continue
		}

		name := ft.Name
		flagName := strings.ToLower(name[:1]) + name[1:]
		f := ccValue.Field(i)
		p := f.Addr().Interface()

		if flagVar(ft, p, flagName, name) {
			usages = append(usages, `	 -`+flagName+" \t\t "+ft.Tag.Get("usage"))
		}
	}

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `Usage: snmp [options] Oids...
	 -m     Method to operate, get/walk/trapsend/translate (default is get/walk)
	 -c     String Default SNMP community (default "public")
	 -t     One or more SNMP targets (eg. -t 192.168.1.1 -t myCommunity@192.168.1.2:1234)
	 -x/y/z One or more x/y/z vars (eg. -x 1-3 -y 1,3,5 -z 1,2-5)
	 -o     One or more Oids
	 -s     Trap server listening address(eg. :9162)
	 -V     Debug/desc, Verbose logging of packets, oid units, oid description and etc.
`+strings.Join(usages, "\n")+"\n")
	}

	flag.Parse()

	o.mib = snmpp.LoadMibs()
	o.Oids = append(o.Oids, flag.Args()...)
	isTranslate := o.Operate == "translate"
	o.Oids = interpolate(isTranslate, o.mib, o.Oids, util.ExpandNums(x), "x")
	o.Oids = interpolate(isTranslate, o.mib, o.Oids, util.ExpandNums(y), "y")
	o.Oids = interpolate(isTranslate, o.mib, o.Oids, util.ExpandNums(z), "z")

	if strings.Contains(o.Verbose, "debug") {
		logger := gosnmp.NewLogger(log.New(log.Writer(), log.Prefix(), log.Flags()))
		o.Logger = &logger
		log.Printf("Oids: %v", o.Oids)
	}
}

func flagVar(ft reflect.StructField, p any, flagName string, name string) bool {
	if ft.Type == timeDurationType {
		flag.DurationVar(p.(*time.Duration), flagName, 0, "")
		return true
	}

	switch ft.Type.Kind() {
	case reflect.Slice:
		switch ft.Type.Elem().Kind() {
		case reflect.String:
			pp := p.(*[]string)
			flag.Var(&ArrayFlags{pp: pp}, name, "")
			return true
		}
	case reflect.String:
		flag.StringVar(p.(*string), flagName, "", "")
		return true
	case reflect.Bool:
		flag.BoolVar(p.(*bool), flagName, false, "")
		return true
	case reflect.Int:
		flag.IntVar(p.(*int), flagName, 0, "")
		return true
	}

	return false
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
