// Package telnet provides a zgrab2 module that scans for telnet daemons.
// Default Port: 23 (TCP)
//
// The --max-read-size flag allows setting a ceiling to the number of bytes
// that will be read for the banner.
//
// The scan negotiates the options and attempts to grab the banner, using the
// same behavior as the original zgrab.
//
// The output contains the banner and the negotiated options, in the same
// format as the original zgrab.
package telnet

import (
	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *module.Telnet
}

// Init initializes the Scanner.
func (scanner *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.Telnet)
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
	return "telnet"
}

// Scan connects to the target (default port TCP 23) and attempts to grab the Telnet banner.
func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	conn, err := target.Open(&scanner.config.Base)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer conn.Close()
	result := new(TelnetLog)
	if err := GetTelnetBanner(result, conn, scanner.config.MaxReadSize); err != nil {
		if scanner.config.Banner && len(result.Banner) > 0 {
			return zgrab2.TryGetScanStatus(err), result, err
		} else {
			return zgrab2.TryGetScanStatus(err), result.getResult(), err
		}
	}
	return zgrab2.SCAN_SUCCESS, result, nil
}
