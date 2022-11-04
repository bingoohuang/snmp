package snmpp

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	g "github.com/gosnmp/gosnmp"
)

type ClientConfig struct {
	Timeout time.Duration `toml:"timeout" usage:"Timeout to wait for a response."`

	Retries              int    `toml:"retries" usage:"Number of retries to attempt."`
	Version              int    `toml:"version" usage:"SNMP version; can be 1, 2, or 3."`
	UnconnectedUDPSocket bool   `toml:"unconnected_udp_socket" usage:"Unconnected UDP socket When true, SNMP responses are accepted from any address not just the requested address. This can be useful when gathering from redundant/failover systems."`
	Community            string `toml:"community" usage:"SNMP community string, Parameters for Version 1 & 2"`

	MaxRepetitions int `toml:"max_repetitions" usage:"The GETBULK max-repetitions parameter, Parameters for Version 2 & 3"`

	// Parameters for Version 3
	// SNMPv3 authentication and encryption options.

	ContextName  string `toml:"context_name" usage:"Context Name. for SNMPv3"`
	SecLevel     string `toml:"sec_level" usage:"Security Level; one of none, authNoPriv, or authPriv. for SNMPv3"`
	UserName     string `toml:"user_name" usage:"User Name. for SNMPv3"`
	AuthProtocol string `toml:"auth_protocol" usage:"Authentication protocol; one of MD5, SHA, SHA224, SHA256, SHA384, SHA512. for SNMPv3"`
	AuthPassword string `toml:"auth_password" usage:"Authentication password"`
	// Protocols "AES192", "AES192", "AES256", and "AES256C" require the underlying net-snmp tools
	// to be compiled with --enable-blumenthal-aes (http://www.net-snmp.org/docs/INSTALL.html)
	PrivProtocol string `toml:"priv_protocol" usage:"Privacy protocol used for encrypted messages; one of DES, AES, AES192, AES192C, AES256, AES256C"`
	PrivPassword string `toml:"priv_password" usage:"Privacy password used for encrypted messages"`

	EngineID    string `toml:"engine_id"`
	EngineBoots int    `toml:"engine_boots"`
	EngineTime  int    `toml:"engine_time"`
}

const (
	DefaultSnmpPort  = 161
	DefaultCommunity = "public"
)

func If[T any](b bool, s1, s2 T) T {
	if b {
		return s1
	}

	return s2
}

func Or(a, b string) string {
	if a == "" {
		return b
	}

	return a
}

func (o *ClientConfig) CreateGoSNMP(target string, verbose bool) *g.GoSNMP {
	gs := &g.GoSNMP{
		Port:               DefaultSnmpPort,
		Transport:          "udp",
		Version:            g.Version2c,
		Retries:            If(o.Retries <= 0, 3, o.Retries),
		MaxRepetitions:     uint32(If(o.MaxRepetitions <= 0, 10, o.MaxRepetitions)),
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

	if o.UserName != "" {
		gs.Version = g.Version3
	}

	if gs.Version == g.Version3 {
		o.SetVersion3Parameters(gs)
	} else {
		gs.Community = Or(o.Community, DefaultCommunity)
	}

	if verbose {
		setupVerbose(gs)
	}

	target = o.parseCommunity(target, gs)
	if err := setTarget(gs, target); err != nil {
		log.Fatalf("set target failed: %v", err)
	}

	return gs
}

func (o *ClientConfig) parseCommunity(target string, gs *g.GoSNMP) string {
	if p := strings.LastIndex(target, "@"); p < 0 {
		gs.Community = o.Community
	} else {
		gs.Community = target[:p]
		target = target[p+1:]
	}

	return target
}

func setupVerbose(gs *g.GoSNMP) {
	// Function handles for collecting metrics on query latencies.
	var sent time.Time
	gs.OnSent = func(*g.GoSNMP) { sent = time.Now() }
	gs.OnRecv = func(*g.GoSNMP) { log.Printf("Query latency: %s", time.Since(sent)) }
}

func (o *ClientConfig) SetVersion3Parameters(gs *g.GoSNMP) {
	if o.ContextName != "none" {
		gs.ContextName = Or(o.ContextName, "public")
	}

	sp := &g.UsmSecurityParameters{}
	gs.SecurityParameters = sp
	gs.SecurityModel = g.UserSecurityModel

	switch strings.ToLower(o.SecLevel) {
	case "none":
		gs.MsgFlags = g.NoAuthNoPriv
	case "authnopriv":
		gs.MsgFlags = g.AuthNoPriv
	case "authpriv":
		gs.MsgFlags = g.AuthPriv
	default:
		if o.AuthPassword != "" && o.PrivPassword != "" {
			gs.MsgFlags = g.AuthPriv
		} else if o.AuthPassword != "" && o.PrivPassword == "" {
			gs.MsgFlags = g.AuthNoPriv
		} else if o.AuthPassword == "" && o.PrivPassword == "" {
			gs.MsgFlags = g.NoAuthNoPriv
		}
	}

	sp.UserName = o.UserName

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
	case "none":
		sp.AuthenticationProtocol = g.NoAuth
	default:
		if o.AuthPassword != "" {
			sp.AuthenticationProtocol = g.MD5
		} else {
			sp.AuthenticationProtocol = g.NoAuth
		}
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
	case "none":
		sp.PrivacyProtocol = g.NoPriv
	default:
		if o.PrivPassword != "" {
			sp.PrivacyProtocol = g.DES
		} else {
			sp.PrivacyProtocol = g.NoPriv
		}
	}

	sp.PrivacyPassphrase = o.PrivPassword
	sp.AuthoritativeEngineID = o.EngineID
	sp.AuthoritativeEngineBoots = uint32(o.EngineBoots)
	sp.AuthoritativeEngineTime = uint32(o.EngineTime)
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
