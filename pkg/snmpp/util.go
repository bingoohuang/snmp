package snmpp

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/bingoohuang/snmp/pkg/smi"
)

func userMibDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Abs(filepath.Join(usr.HomeDir, ".snmp", "mibs"))
}

func LoadMibs() *smi.MIB {
	var dirs []string

	if dir, err := userMibDir(); err != nil {
		log.Printf("W! failed to find user mib dir(~/.snmp/mibs), error: %v", err)
	} else {
		dirs = append(dirs, dir)
	}

	const sysMibDir = "/usr/share/snmp/mibs"
	if d, err := os.Stat(sysMibDir); err == nil && d.IsDir() {
		dirs = append(dirs, sysMibDir)
	}

	mib := smi.NewMIB(dirs...)
	if err := mib.LoadModules(); err != nil {
		log.Printf("W! failed to load mibs, error: %v", err)
	}

	return mib
}

func ParseOIDSymbolName(dotOid string, mib *smi.MIB) string {
	oid, err := smi.ParseOID(dotOid)
	if err != nil {
		log.Printf("E! parse oid error %v", err)
		return dotOid
	}

	if symbol, suffix := mib.Symbol(oid); symbol != nil {
		return SymbolString(symbol, suffix)
	}

	return "Unknown"
}

func IsSymbolName(oid string) bool {
	return strings.Contains(oid, "::")
}

func SymbolString(symbol *smi.Symbol, suffix smi.OID) string {
	if len(suffix) == 0 {
		return symbol.String()
	}

	return fmt.Sprintf("%s.%s", symbol, suffix.String())
}
