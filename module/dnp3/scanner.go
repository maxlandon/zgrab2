// Package dnp3 provides a zgrab2 module that scans for dnp3.
// Default port: 20000 (TCP)
//
// Copied unmodified from the original zgrab.
// Connects, and reads the banner. Returns the raw response.
package dnp3

import (
	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *module.DNP3
}

// Init initializes the Scanner.
func (scanner *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.DNP3)
	scanner.config = f
	return nil
}

// InitPerSender initializes the scanner for a given sender.
func (scanner *Scanner) InitPerSender(senderID int) error {
	return nil
}

// GetName returns the Scanner name defined in the Flags.
func (scanner *Scanner) GetName() string {
	return scanner.config.Name
}

// GetTrigger returns the Trigger defined in the Flags.
func (scanner *Scanner) GetTrigger() string {
	return scanner.config.Trigger
}

// Protocol returns the protocol identifier of the scan.
func (scanner *Scanner) Protocol() string {
	return "dnp3"
}

// Scan probes for a DNP3 service.
// Connects to the configured TCP port (default 20000) and reads the banner.
func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	// TODO: Allow UDP?
	conn, err := target.Open(&scanner.config.Base)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer conn.Close()
	ret := new(DNP3Log)
	if err := GetDNP3Banner(ret, conn); err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	return zgrab2.SCAN_SUCCESS, ret, nil
}
