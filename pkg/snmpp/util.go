package snmpp

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
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

func LoadMibs(verbose bool) *smi.MIB {
	var dirs []string

	dirs = addDir("/usr/share/snmp/mibs", dirs, verbose)

	userDir, err := userMibDir()
	if err != nil {
		if verbose {
			log.Printf("W! failed to read %s, error: %v", "~/.snmp/mibs", err)
		}
	} else {
		dirs = addDir(userDir, dirs, false)
	}

	mib := smi.NewMIB(dirs...)
	if err := mib.LoadModules(); err != nil {
		log.Printf("W! failed to load mibs, error: %v", err)
	}

	return mib
}

func addDir(sysMibDir string, dirs []string, verbose bool) []string {
	d, err := os.Stat(sysMibDir)
	switch {
	case err != nil:
		if verbose {
			log.Printf("W! failed to read %s, error: %v", sysMibDir, err)
		}
	case !d.IsDir():
		if verbose {
			log.Printf("W! %s is not a directory", sysMibDir)
		}
	default:
		dirs = append(dirs, sysMibDir)
	}

	return dirs
}

func ParseOIDSymbolName(dotOid string, mib *smi.MIB) (symbolName, description string, sym *smi.Symbol) {
	oid, err := smi.ParseOID(dotOid)
	if err != nil {
		log.Printf("E! parse oid error %v", err)
		return "", "", nil
	}

	if symbol, suffix := mib.Symbol(oid); symbol != nil {
		symbolString, desc := SymbolString(symbol, suffix)
		return symbolString, desc, symbol
	}

	return "Unknown", "", nil
}

var OidPattern = regexp.MustCompile(`\.?\d+(\.\d+)*`)

func IsSymbolName(oid string) bool {
	return strings.Contains(oid, "::") || !OidPattern.MatchString(oid)
}

func JoinLines(s string) string {
	return regexp.MustCompile(`\s\s+`).ReplaceAllString(s, " ")
}

func SymbolString(symbol *smi.Symbol, suffix smi.OID) (symbolName, description string) {
	s := symbol.String()
	if len(suffix) == 0 {
		return s, JoinLines(symbol.Description)
	}

	return fmt.Sprintf("%s.%s", s, suffix.String()), JoinLines(symbol.Description)
}
