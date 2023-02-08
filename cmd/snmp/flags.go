package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/bingoohuang/gg/pkg/v"
	"github.com/bingoohuang/snmp/pkg/smi"
	"github.com/bingoohuang/snmp/pkg/snmpp"
	"github.com/bingoohuang/snmp/pkg/util"
	"github.com/gosnmp/gosnmp"
	"github.com/spf13/pflag"
)

type Options struct {
	snmpp.ClientConfig

	Agents  []string
	Oids    []string
	Verbose string

	TrapAddr string
	Operate  string
	Logger   *gosnmp.Logger
	mib      *smi.MIB

	Version bool
}

var timeDurationType = reflect.TypeOf(time.Duration(0))

func (o *Options) ParseFlags() {
	var x, y, z []string

	pflag.StringVarP(&o.Operate, "method", "m", "get/walk", "")
	pflag.StringArrayVarP(&o.Agents, "agent", "t", nil, `Agent addresses to retrieve values from.`+"\n"+
		`   format:  agents = <community@><scheme://><hostname>:<optional port>`+"\n"+
		`   scheme:  optional, either udp, udp4, udp6, tcp, tcp4, tcp6. default is udp`+"\n"+
		`   example: 127.0.0.1, myCommunity@192.168.1.2:1234, udp://127.0.0.1:161, tcp://127.0.0.1:161, udp4://v4only-snmp-agent`+"\n")
	pflag.StringArrayVarP(&x, "vx", "x", nil, "x vars (eg. -x 1-3 -x 5)")
	pflag.StringArrayVarP(&y, "vy", "y", nil, "y var")
	pflag.StringArrayVarP(&z, "vz", "z", nil, "z var")
	pflag.StringArrayVarP(&o.Oids, "oid", "o", nil, "oids")
	pflag.StringVarP(&o.TrapAddr, "trapAddr", "", "", "Trap server listening address(eg. :9162)")
	pflag.StringVarP(&o.Verbose, "verbose", "V", "", "debug/desc, Verbose logging of packets, oid units, oid description and etc.")
	pflag.BoolVarP(&o.Version, "ver", "", false, "print snmp version and exit")

	cnf := &o.ClientConfig
	declarePflags(pflag.CommandLine, cnf)

	pflag.Parse()

	if o.Version {
		fmt.Println(v.Version())
		os.Exit(0)
	}

	debug := strings.Contains(o.Verbose, "debug")
	o.mib = snmpp.LoadMibs(debug)
	o.Oids = append(o.Oids, pflag.Args()...)
	isTranslate := o.Operate == "translate"
	o.Oids = interpolate(isTranslate, o.mib, o.Oids, util.ExpandNums(x), "x")
	o.Oids = interpolate(isTranslate, o.mib, o.Oids, util.ExpandNums(y), "y")
	o.Oids = interpolate(isTranslate, o.mib, o.Oids, util.ExpandNums(z), "z")
	o.Oids = unique(o.Oids)

	if debug {
		logger := gosnmp.NewLogger(log.New(log.Writer(), log.Prefix(), log.Flags()))
		o.Logger = &logger
		log.Printf("Oids: %v", o.Oids)
	}
}

func declarePflags(pf *pflag.FlagSet, cnf any) {
	ccValue := reflect.ValueOf(cnf).Elem()
	ccType := ccValue.Type()
	var usages []string
	for i := 0; i < ccType.NumField(); i++ {
		ft := ccType.Field(i)
		if !ft.IsExported() {
			continue
		}

		tagName := ft.Tag.Get("flag")
		if tagName == "-" {
			continue
		}

		name := ft.Name
		flagName := strings.ToLower(name[:1]) + name[1:]
		f := ccValue.Field(i)
		p := f.Addr().Interface()

		if flagVar(pf, ft, p, flagName, name) {
			usages = append(usages, `	 -`+flagName+" \t\t "+ft.Tag.Get("usage"))
		}
	}
}

func unique[T comparable](src []T) []T {
	keys := make(map[T]bool)
	list := make([]T, 0, len(src))
	for _, entry := range src {
		if _, ok := keys[entry]; !ok {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func flagVar(pf *pflag.FlagSet, ft reflect.StructField, p any, flagName string, name string) bool {
	if ft.Type == timeDurationType {
		pf.DurationVar(p.(*time.Duration), flagName, 0, "")
		return true
	}

	usage := ft.Tag.Get("usage")
	value := ft.Tag.Get("value")
	tagName := ft.Tag.Get("flag")
	shortName := ""
	if tagName != "" {
		if i := strings.Index(tagName, ","); i >= 0 {
			shortName = tagName[i+1:]
			tagName = tagName[:i]
		}
	}

	if tagName != "" {
		name = tagName
	}

	switch ft.Type.Kind() {
	case reflect.Slice:
		switch ft.Type.Elem().Kind() {
		case reflect.String:
			pf.StringArrayVarP(p.(*[]string), name, shortName, strings.Split(value, ","), usage)
			return true
		}
	case reflect.String:
		pf.StringVarP(p.(*string), flagName, shortName, value, usage)
		return true
	case reflect.Bool:
		pf.BoolVarP(p.(*bool), flagName, shortName, value == "true", usage)
		return true
	case reflect.Int:
		intValue, _ := strconv.Atoi(value)
		pf.IntVarP(p.(*int), flagName, shortName, intValue, usage)
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
			continue
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
