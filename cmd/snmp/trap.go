package main

import (
	"log"
	"net"

	g "github.com/gosnmp/gosnmp"
)

func (o *Options) trap() {
	if o.TrapAddr == "" {
		return
	}

	tl := g.NewTrapListener()
	tl.OnNewTrap = o.trapHandler
	tl.Params = g.Default
	if o.Logger != nil {
		tl.Params.Logger = *o.Logger
	}

	if err := tl.Listen(o.TrapAddr); err != nil {
		log.Printf("E! error in listen: %s", err)
	}
}

func (o *Options) trapHandler(p *g.SnmpPacket, addr *net.UDPAddr) {
	log.Printf("got trapdata from %s", addr.IP)
	for i, v := range p.Variables {
		o.printPdu("trap", addr.String(), i, v, 0)
	}
}
