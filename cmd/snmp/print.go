package main

import (
	"fmt"

	"github.com/bingoohuang/snmp/pkg/snmpp"
	g "github.com/gosnmp/gosnmp"
)

func (o *Options) printPdu(typ, target string, i int, pdu g.SnmpPDU) {
	symbol, description := snmpp.ParseOIDSymbolName(pdu.Name, o.mib)

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
		fmt.Printf("%sName%s: %s\n", KeyStyle, EndStyle, symbol)
		fmt.Printf("%sDescription%s: %s\n", KeyStyle, EndStyle, description)
		symbol = ""
	}

	if symbol != "" {
		symbol = "[" + symbol + "]"
	}

	arrow := StringStyle + "=>" + EndStyle

	fmt.Printf("%s%s[%d]%s[%s] %s %v: ", typ, target, i, symbol, pdu.Name, arrow, pdu.Type)

	switch pdu.Type {
	case g.OctetString:
		fmt.Printf("%s\n", pdu.Value.([]byte))
	case g.ObjectIdentifier:
		fmt.Printf("%s\n", pdu.Value.(string))
	default:
		fmt.Printf("%v\n", pdu.Value)
	}
}
