package main

import (
	"log"
	"strings"

	g "github.com/gosnmp/gosnmp"
)

func (t *Target) snmpWalk() {
	if !strings.Contains(t.Operate, "walk") {
		return
	}

	for _, oid := range t.Oids {
		i := 0
		if err := t.BulkWalk(oid, func(pdu g.SnmpPDU) error {
			t.printPdu("walk", t.target, i, pdu)
			i++
			return nil
		}); err != nil {
			log.Printf("W! snmpwalk error: %v", err)
		}
	}
}

func (t *Target) snmpGet() {
	if !strings.Contains(t.Operate, "get") {
		return
	}

	result, err := t.Get(t.Oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		log.Printf("W! snmpget %v error: %v", t.Oids, err)
		return
	}

	for i, pdu := range result.Variables {
		t.printPdu("get", t.target, i, pdu)
	}
}
