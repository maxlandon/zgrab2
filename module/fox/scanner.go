// Package fox provides a zgrab2 module that scans for fox.
// Default port: 1911 (TCP)
//
// Copied unmodified from the original zgrab.
// Connects, sends a static query, and reads the banner. Parses out as much of the response as possible.
package fox

import (
	"errors"

	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *module.Fox
}

// Init initializes the Scanner.
func (scanner *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.Fox)
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
	return "fox"
}

// Scan probes for a Tridium Fox service.
// 1. Opens a TCP connection to the configured port (default 1911)
// 2. Sends a static query
// 3. Attempt to read the response (up to 8k + 4 bytes -- larger responses trigger an error)
// 4. If the response has the Fox response prefix, mark the scan as having detected the service.
// 5. Attempt to read any / all of the data fields from the Log struct.
func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	conn, err := target.Open(&scanner.config.Base)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer conn.Close()
	result := new(FoxLog)

	err = GetFoxBanner(result, conn)
	if !result.IsFox {
		result = nil
		err = &zgrab2.ScanError{
			Err:    errors.New("host responds, but is not a fox service"),
			Status: zgrab2.SCAN_PROTOCOL_ERROR,
		}
	}
	return zgrab2.TryGetScanStatus(err), result, err
}
