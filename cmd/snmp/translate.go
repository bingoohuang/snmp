package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bingoohuang/snmp/pkg/smi"
	"github.com/bingoohuang/snmp/pkg/snmpp"
)

func (o *Options) translate() {
	if !strings.Contains(o.Mode, "translate") {
		return
	}

	for _, v := range o.Oids {
		if snmpp.IsSymbolName(v) {
			if oid, err := o.mib.OID(v); err != nil {
				fmt.Printf("%s => error %v\n", v, err)
			} else {
				fmt.Printf("%s => %s\n", v, oid)
			}
		} else {
			oid, err := smi.ParseOID(v)
			if err != nil {
				fmt.Printf("%s => error %v\n", v, err)
				continue
			}

			if symbol, suffix := o.mib.Symbol(oid); symbol == nil {
				fmt.Printf("%s => unknown\n", v)
			} else {
				fmt.Printf("%s => %s\n", v, snmpp.SymbolString(symbol, suffix))
			}
		}
	}

	os.Exit(0)
}
