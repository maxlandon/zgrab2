// Ref: https://github.com/salesforce/jarm
// https://engineering.salesforce.com/easily-identify-malicious-servers-on-the-internet-with-jarm-e095edac525a?gi=4dd05e2277e4
package jarm

import (
	_ "fmt"
	"net"
	"strings"
	"time"

	jarm "github.com/RumbleDiscovery/jarm-go"
	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// Scanner is the implementation of the zgrab2.Scanner interface.
type Scanner struct {
	config *module.JARM
}

type Results struct {
	Fingerprint string `json:"fingerprint"`
	error       string `json:"error,omitempty"`
}

// GetName returns the Scanner name defined in the Flags.
func (scanner *Scanner) GetName() string {
	return scanner.config.Name
}

// GetPort returns the port being scanned.
func (scanner *Scanner) GetPort() uint {
	return scanner.config.Port
}

// GetTrigger returns the Trigger defined in the Flags.
func (scanner *Scanner) GetTrigger() string {
	return scanner.config.Trigger
}

// Protocol returns the protocol identifier of the scan.
func (scanner *Scanner) Protocol() string {
	return "jarm"
}

// InitPerSender initializes the scanner for a given sender.
func (scanner *Scanner) InitPerSender(senderID int) error {
	return nil
}

// Init initializes the Scanner with the command-line flags.
func (scanner *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.JARM)
	scanner.config = f
	return nil
}

func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	// Stores raw hashes returned from parsing each protocols Hello message
	rawhashes := []string{}

	// Loop through each Probe type
	for _, probe := range jarm.GetProbes(target.Host(), int(scanner.GetPort())) {
		var (
			conn net.Conn
			err  error
			ret  []byte
		)
		conn, err = target.Open(&scanner.config.Base)
		if err != nil {
			return zgrab2.TryGetScanStatus(err), nil, err
		}

		_, err = conn.Write(jarm.BuildProbe(probe))
		if err != nil {
			rawhashes = append(rawhashes, "")
			conn.Close()
			continue
		}

		ret, _ = zgrab2.ReadAvailableWithOptions(conn, 1484, 500*time.Millisecond, 0, 1484)

		ans, err := jarm.ParseServerHello(ret, probe)
		if err != nil {
			rawhashes = append(rawhashes, "")
			conn.Close()
			continue
		}

		rawhashes = append(rawhashes, ans)
		conn.Close()
	}

	return zgrab2.SCAN_SUCCESS, &Results{
		Fingerprint: jarm.RawHashToFuzzyHash(strings.Join(rawhashes, ",")),
	}, nil
}
