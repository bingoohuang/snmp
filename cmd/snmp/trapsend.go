package main

import (
	"log"
	"os"
	"strings"

	g "github.com/gosnmp/gosnmp"
)

func (t *Target) trapSend() {
	if !strings.Contains(t.Operate, "trapsend") {
		return
	}

	trap := g.SnmpTrap{
		Variables: []g.SnmpPDU{{
			Name:  "1.3.6.1.2.1.1.6",
			Type:  g.ObjectIdentifier,
			Value: "1.3.6.1.2.1.1.6.10",
		}, {
			Name:  "1.3.6.1.2.1.1.7",
			Type:  g.OctetString,
			Value: "Testing TCP trap...",
		}, {
			Name:  "1.3.6.1.2.1.1.8",
			Type:  g.Integer,
			Value: 123,
		}},
	}

	if _, err := t.SendTrap(trap); err != nil {
		log.Printf("E! SendTrap() err: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
