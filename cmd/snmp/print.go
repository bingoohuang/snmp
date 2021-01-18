package main

import (
	"fmt"

	"github.com/bingoohuang/snmp/pkg/snmpp"
	g "github.com/gosnmp/gosnmp"
)

func (o *Options) printPdu(typ, target string, i int, pdu g.SnmpPDU) {
	symbolName, description, syn := snmpp.ParseOIDSymbolName(pdu.Name, o.mib)

	if o.Operate != typ {
		typ = "[" + typ + "]"
	} else {
		typ = ""
	}

	if len(o.Targets) > 1 {
		target = "[" + target + "]"
	} else {
		target = ""
	}

	if symbolName != "" {
		symbolName = "[" + KeyStyle + symbolName + EndStyle + "]"
	}

	arrow := StringStyle + "=>" + EndStyle

	fmt.Printf("%s%s[%d]%s[%s] %s %v: ", typ, target, i, symbolName, pdu.Name, arrow, pdu.Type)
	fmt.Print(RedStyle)
	switch pdu.Type {
	case g.OctetString:
		fmt.Printf("%s", pdu.Value.([]byte))
	case g.ObjectIdentifier:
		fmt.Printf("%s", pdu.Value.(string))
	default:
		fmt.Printf("%v", pdu.Value)
	}
	fmt.Print(EndStyle)

	if o.Verbose && description != "" {
		if syn.Unit != "" {
			fmt.Printf(" Unit: %s", KeyStyle+syn.Unit+EndStyle)
		}

		fmt.Printf(" Desc: %s", GrayStyle+description+EndStyle)
	}

	fmt.Println()
}
