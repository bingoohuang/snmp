package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bingoohuang/gg/pkg/ss"
	g "github.com/gosnmp/gosnmp"
)

func main() {
	options := Options{}
	options.ParseFlags()

	for _, t := range options.Agents {
		options.do(t)
	}

	options.translate()
	options.trap()
}

type Target struct {
	*g.GoSNMP
	*Options
	target string
}

func (t *Target) Close() {
	_ = t.Conn.Close()
}

func (o *Options) do(target string) {
	t := o.createAgent(target)
	if err := t.Connect(); err != nil {
		log.Printf("E! Connect() err: %v", err)
		os.Exit(1)
	}

	defer t.Close()

	t.trapSend()
	t.snmpGet()
	t.snmpWalk()
}

const (
	DefaultSnmpPort  = 161
	DefaultCommunity = "public"
)

func (o *Options) createAgent(target string) Target {
	gs := &g.GoSNMP{
		Port:               DefaultSnmpPort,
		Transport:          "udp",
		Version:            g.Version2c,
		Retries:            ss.Ifi(o.Retries <= 0, 3, o.Retries),
		MaxRepetitions:     uint32(ss.Ifi(o.MaxRepetitions <= 0, 10, o.MaxRepetitions)),
		Timeout:            time.Duration(10) * time.Second,
		ExponentialTimeout: false,
		MaxOids:            g.MaxOids,
	}

	if o.Timeout > 0 {
		gs.Timeout = o.Timeout
	}

	switch o.Version {
	case 3:
		gs.Version = g.Version3
	case 1:
		gs.Version = g.Version1
	default:
		gs.Version = g.Version2c
	}

	if o.Version == 3 {
		o.setVersion3Parameters(gs)
	} else {
		gs.Community = ss.Or(o.Community, DefaultCommunity)
	}

	o.setupVerbose(gs)
	target = o.parseCommunity(target, gs)
	if err := setTarget(gs, target); err != nil {
		log.Fatalf("set target failed: %v", err)
	}

	return Target{
		GoSNMP:  gs,
		target:  refineTargetForOutput(gs),
		Options: o,
	}
}

func (o *Options) setVersion3Parameters(gs *g.GoSNMP) {
	gs.ContextName = o.ContextName

	sp := &g.UsmSecurityParameters{}
	gs.SecurityParameters = sp
	gs.SecurityModel = g.UserSecurityModel

	switch strings.ToLower(o.SecLevel) {
	case "noauthnopriv", "":
		gs.MsgFlags = g.NoAuthNoPriv
	case "authnopriv":
		gs.MsgFlags = g.AuthNoPriv
	case "authpriv":
		gs.MsgFlags = g.AuthPriv
	default:
		gs.MsgFlags = g.NoAuthNoPriv
	}

	sp.UserName = o.SecName

	switch strings.ToLower(o.AuthProtocol) {
	case "md5":
		sp.AuthenticationProtocol = g.MD5
	case "sha":
		sp.AuthenticationProtocol = g.SHA
	case "sha224":
		sp.AuthenticationProtocol = g.SHA224
	case "sha256":
		sp.AuthenticationProtocol = g.SHA256
	case "sha384":
		sp.AuthenticationProtocol = g.SHA384
	case "sha512":
		sp.AuthenticationProtocol = g.SHA512
	case "":
		sp.AuthenticationProtocol = g.NoAuth
	default:
		sp.AuthenticationProtocol = g.NoAuth
	}

	sp.AuthenticationPassphrase = o.AuthPassword

	switch strings.ToLower(o.PrivProtocol) {
	case "des":
		sp.PrivacyProtocol = g.DES
	case "aes":
		sp.PrivacyProtocol = g.AES
	case "aes192":
		sp.PrivacyProtocol = g.AES192
	case "aes192c":
		sp.PrivacyProtocol = g.AES192C
	case "aes256":
		sp.PrivacyProtocol = g.AES256
	case "aes256c":
		sp.PrivacyProtocol = g.AES256C
	case "":
		sp.PrivacyProtocol = g.NoPriv
	default:
		sp.PrivacyProtocol = g.NoPriv
	}

	sp.PrivacyPassphrase = o.PrivPassword
	sp.AuthoritativeEngineID = o.EngineID
	sp.AuthoritativeEngineBoots = uint32(o.EngineBoots)
	sp.AuthoritativeEngineTime = uint32(o.EngineTime)
}

func (o *Options) setupVerbose(gs *g.GoSNMP) {
	if o.Logger == nil {
		return
	}

	gs.Logger = *o.Logger

	// Function handles for collecting metrics on query latencies.
	var sent time.Time
	gs.OnSent = func(*g.GoSNMP) { sent = time.Now() }
	gs.OnRecv = func(*g.GoSNMP) { log.Printf("Query latency: %s", time.Since(sent)) }
}

func (o *Options) parseCommunity(target string, gs *g.GoSNMP) string {
	if p := strings.LastIndex(target, "@"); p < 0 {
		gs.Community = o.Community
	} else {
		gs.Community = target[:p]
		target = target[p+1:]
	}

	return target
}

// setTarget takes a url (scheme://host:port) and sets the GoSNMP struct's corresponding fields.
// This shouldn't be called after using the wrapped GoSNMP struct, for example after connecting.
func setTarget(gs *g.GoSNMP, target string) error {
	if !strings.Contains(target, "://") {
		target = "udp://" + target
	}

	u, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("parse target %s failed: %v", target, err)
	}

	// Only allow udp{4,6} and tcp{4,6}.
	// Allowing ip{4,6} does not make sense as specifying a port
	// requires the specification of a protocol.
	// gosnmp does not handle these errors well, which is why
	// they can result in cryptic errors by net.Dial.
	switch u.Scheme {
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
		gs.Transport = u.Scheme
	default:
		return fmt.Errorf("unsupported scheme: %v", u.Scheme)
	}

	gs.Target = u.Hostname()

	portStr := u.Port()
	if portStr == "" {
		portStr = "161"
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return fmt.Errorf("parsing port: %w", err)
	}
	gs.Port = uint16(port)
	return nil
}

func refineTargetForOutput(gs *g.GoSNMP) string {
	target := ""

	if gs.Community != DefaultCommunity {
		target = gs.Community + "@"
	}

	target += gs.Target
	if gs.Port != DefaultSnmpPort {
		target += fmt.Sprintf(":%d", gs.Port)
	}
	return target
}
