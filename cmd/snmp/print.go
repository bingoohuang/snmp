package main

import (
	"fmt"

	"github.com/bingoohuang/snmp/pkg/snmpp"
	g "github.com/gosnmp/gosnmp"
)

func (o *Options) printPdu(typ, target string, i int, pdu g.SnmpPDU) {
	symbol := snmpp.ParseOIDSymbolName(pdu.Name, o.mib)

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

	fmt.Printf("%s%s[%d][%s][%s] = %v: ", typ, target, i, symbol, pdu.Name, pdu.Type)

	switch pdu.Type {
	case g.OctetString:
		fmt.Printf("%s\n", pdu.Value.([]byte))
	case g.ObjectIdentifier:
		fmt.Printf("%s\n", pdu.Value.(string))
	default:
		fmt.Printf("%v\n", pdu.Value)
	}
}
