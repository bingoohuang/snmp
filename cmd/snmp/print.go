package main

import (
	"fmt"

	"github.com/bingoohuang/snmp/pkg/snmpp"
	g "github.com/gosnmp/gosnmp"
)

func (o *Options) printPdu(typ, target string, i int, pdu g.SnmpPDU) {
	symbolName, description, syn := snmpp.ParseOIDSymbolName(pdu.Name, o.mib)

	if o.Mode != typ {
		typ = "[" + typ + "]"
	} else {
		typ = ""
	}

	if len(o.Targets) > 1 {
		target = "[" + target + "]"
	} else {
		target = ""
	}

	if o.Verbose && description != "" {
		fmt.Printf("%sOidName%s: %s", KeyStyle, EndStyle, symbolName)
		if syn.Unit != "" {
			fmt.Printf(" %sUnit%s: %s\n", KeyStyle, EndStyle, syn.Unit)
		} else {
			fmt.Println()
		}
		fmt.Printf("%sDescription%s: %s\n", KeyStyle, EndStyle, description)
		symbolName = ""
	}

	if symbolName != "" {
		symbolName = "[" + symbolName + "]"
	}

	arrow := StringStyle + "=>" + EndStyle

	fmt.Printf("%s%s[%d]%s[%s] %s %v: ", typ, target, i, symbolName, pdu.Name, arrow, pdu.Type)

	switch pdu.Type {
	case g.OctetString:
		fmt.Printf("%s\n", pdu.Value.([]byte))
	case g.ObjectIdentifier:
		fmt.Printf("%s\n", pdu.Value.(string))
	default:
		fmt.Printf("%v\n", pdu.Value)
	}
}
