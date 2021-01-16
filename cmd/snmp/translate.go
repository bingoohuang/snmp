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

	arrow := StringStyle + "=>" + EndStyle
	for _, v := range o.Oids {
		if snmpp.IsSymbolName(v) {
			if oid, err := o.mib.OID(v); err != nil {
				fmt.Printf("%s %s error %v\n", v, arrow, err)
			} else {
				if symbol, suffix := o.mib.Symbol(oid); symbol != nil {
					o.printSymbolWithDescription(symbol, suffix)
				}

				fmt.Printf("%s %s %s\n", v, arrow, oid)
			}
		} else {
			oid, err := smi.ParseOID(v)
			if err != nil {
				fmt.Printf("%s %s error %v\n", v, arrow, err)
				continue
			}

			if symbol, suffix := o.mib.Symbol(oid); symbol == nil {
				fmt.Printf("%s %s unknown\n", v, arrow)
			} else {
				symbolName := o.printSymbolWithDescription(symbol, suffix)

				fmt.Printf("%s %s %s\n", v, arrow, symbolName)
			}
		}
	}

	os.Exit(0)
}

const (
	KeyStyle    = "\x1B[94m"
	StringStyle = "\x1B[92m"
	EndStyle    = "\x1B[0m"
)

func (o *Options) printSymbolWithDescription(symbol *smi.Symbol, suffix smi.OID) string {
	symbolName, description := snmpp.SymbolString(symbol, suffix)
	if o.Verbose && description != "" {
		fmt.Printf("%sObjectType%s: %s", KeyStyle, EndStyle, symbolName)
		if symbol.Unit != "" {
			fmt.Printf(" %sUnit%s: %s\n", KeyStyle, EndStyle, symbol.Unit)
		} else {
			fmt.Println()
		}

		fmt.Printf("%sDescription%s: %s\n", KeyStyle, EndStyle, description)
	}

	return symbolName
}
