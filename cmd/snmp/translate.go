package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bingoohuang/snmp/pkg/smi"
	"github.com/bingoohuang/snmp/pkg/snmpp"
)

func (o *Options) translate() {
	if !strings.Contains(o.Operate, "translate") {
		return
	}

	arrow := StringStyle + "=>" + EndStyle
	for _, v := range o.Oids {
		if snmpp.IsSymbolName(v) {
			if oid, err := o.mib.OID(v); err != nil {
				fmt.Printf("%s %s error %v", v, arrow, err)
			} else {
				fmt.Printf("%s %s %s", v, arrow, oid)
				if symbol, suffix := o.mib.Symbol(oid); symbol != nil {
					o.printSymbolWithDescription(symbol, suffix, false)
				}
			}
		} else {
			oid, err := smi.ParseOID(v)
			if err != nil {
				fmt.Printf("%s %s error %v", v, arrow, err)
				continue
			}

			if symbol, suffix := o.mib.Symbol(oid); symbol == nil {
				fmt.Printf("%s %s unknown", v, arrow)
			} else {
				symbolName, f := o.printSymbolWithDescription(symbol, suffix, true)
				fmt.Printf("%s %s %s", v, arrow, symbolName)
				f()
			}
		}

		fmt.Println()
	}

	os.Exit(0)
}

const (
	KeyStyle    = "\x1B[94m"
	StringStyle = "\x1B[92m"
	RedStyle    = "\x1B[31m"
	GrayStyle   = "\x1B[37m"
	EndStyle    = "\x1B[0m"
)

func (o *Options) printSymbolWithDescription(symbol *smi.Symbol, suffix smi.OID, delay bool) (string, func()) {
	symbolName, description := snmpp.SymbolString(symbol, suffix)
	f := func() {}

	if o.Verbose && description != "" {
		f = func() {
			if symbol.Unit != "" {
				fmt.Printf(" Unit: %s", KeyStyle+symbol.Unit+EndStyle)
			}

			fmt.Printf(" Desc: %s", GrayStyle+description+EndStyle)
		}
	}

	if !delay {
		f()
	}

	return symbolName, f
}
