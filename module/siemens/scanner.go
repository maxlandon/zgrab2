// Package siemens provides a zgrab2 module that scans for Siemens S7.
// Default port: TCP 102
// Ported from the original zgrab. Input and output are identical.
package siemens

import (
	"net"

	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *module.Siemens
}

// Init initializes the Scanner.
func (scanner *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.Siemens)
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
	return "siemens"
}

// Scan probes for Siemens S7 services.
// 1. Connect to TCP port 102
// 2. Send a COTP connection packet with destination TSAP 0x0102, source TSAP 0x0100
// 3. If that fails, reconnect and send a COTP connection packet with destination TSAP 0x0200, source 0x0100
// 4. Negotiate S7
// 5. Request to read the module identification (and store it in the output)
// 6. Request to read the component identification (and store it in the output)
// 7. Return the output.
func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	conn, err := target.Open(&scanner.config.Base)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer conn.Close()
	result := new(S7Log)

	err = GetS7Banner(result, conn, func() (net.Conn, error) { return target.Open(&scanner.config.Base) })
	if !result.IsS7 {
		result = nil
	}
	return zgrab2.TryGetScanStatus(err), result, err
}
