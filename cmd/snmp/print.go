package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bingoohuang/snmp/pkg/snmpp"
	g "github.com/gosnmp/gosnmp"
)

func (o *Options) printPdu(typ, target string, i int, pdu g.SnmpPDU, cost time.Duration) {
	symbolName, description, syn := snmpp.ParseOIDSymbolName(pdu.Name, o.mib)

	if o.Operate != typ {
		typ = "[" + typ + "]"
	} else {
		typ = ""
	}

	if len(o.Agents) > 1 {
		target = "[" + target + "]"
	} else {
		target = ""
	}

	if symbolName != "" {
		symbolName = "[" + KeyStyle + symbolName + EndStyle + "]"
	}

	fmt.Printf("%s%s[%d]%s[%s] %s=>%s %v: %s%s%s cost: %s",
		typ, target, i, symbolName, pdu.Name, StringStyle, EndStyle,
		pdu.Type, RedStyle, pduValue(pdu), EndStyle, cost)

	if strings.Contains(o.Verbose, "desc") && description != "" {
		if syn.Unit != "" {
			fmt.Printf(" Unit: %s", KeyStyle+syn.Unit+EndStyle)
		}

		fmt.Printf(" Desc: %s", GrayStyle+description+EndStyle)
	}

	fmt.Println()
}

func pduValue(pdu g.SnmpPDU) any {
	switch pdu.Type {
	case g.OctetString:
		return pdu.Value.([]byte)
	case g.ObjectIdentifier:
		return pdu.Value.(string)
	default:
		return fmt.Sprintf("%v", pdu.Value)
	}
}
