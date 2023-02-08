package main

import (
	"log"
	"strings"
	"time"

	g "github.com/gosnmp/gosnmp"
)

func (t *Target) snmpWalk(pduNames map[string]bool) {
	if !strings.Contains(t.Operate, "walk") {
		return
	}

	for _, oid := range t.Oids {
		i := 0
		start := time.Now()
		if err := t.BulkWalk(oid, func(pdu g.SnmpPDU) error {
			if !pduNames[pdu.Name] {
				pduNames[pdu.Name] = true
				t.printPdu("walk", t.target, i, pdu, time.Since(start))
			}
			i++
			return nil
		}); err != nil {
			log.Printf("W! snmpwalk error: %v", err)
		}
	}
}

func (t *Target) snmpGet(pduNames map[string]bool) {
	if !strings.Contains(t.Operate, "get") {
		return
	}

	start := time.Now()
	result, err := t.Get(t.Oids) // Get() accepts up to g.MAX_OIDS
	cost := time.Since(start)
	if err != nil {
		log.Printf("W! snmpget %v error: %v", t.Oids, err)
		return
	}

	for i, pdu := range result.Variables {
		if !pduNames[pdu.Name] {
			pduNames[pdu.Name] = true
			pduNames[pdu.Name] = true
			t.printPdu("get", t.target, i, pdu, cost)
		}
	}
}
