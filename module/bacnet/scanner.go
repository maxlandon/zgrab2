// Package bacnet provides a zgrab2 module that scans for bacnet.
// Default Port: 47808 / 0xBAC0 (UDP)
//
// Behavior and output copied identically from original zgrab.
package bacnet

import (
	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// Scan results are in log.go

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *module.Bacnet
}

// Init initializes the Scanner.
func (scanner *Scanner) Init(data module.Scan) error {
	f, _ := data.(*module.Bacnet)
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
	return "bacnet"
}

// Scan probes for a BACNet service.
// Behavior taken from original zgrab.
// Connects to the configured port over UDP (default 47808/0xBAC0).
// Attempts to query the following in sequence; if any fails, returning anything that has been detected so far.
// (Unless QueryDeviceID fails, the service is considered to be detected)
// 1. Device ID
// 2. Vendor Number
// 3. Vendor Name
// 4. Firmware Revision
// 5. App software revision
// 6. Object name
// 7. Model  name
// 8. Description
// 9. Location
// The result is a bacnet.Log, and contains any of the above.
func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	conn, err := target.OpenUDP(&scanner.config.Base, &scanner.config.UDP)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer conn.Close()
	ret := new(Log)
	// TODO: if one fails, try others?
	// TODO: distinguish protocol vs app errors
	if err := ret.QueryDeviceID(conn); err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	if err := ret.QueryVendorNumber(conn); err != nil {
		return zgrab2.TryGetScanStatus(err), ret, nil
	}
	if err := ret.QueryVendorName(conn); err != nil {
		return zgrab2.TryGetScanStatus(err), ret, nil
	}
	if err := ret.QueryFirmwareRevision(conn); err != nil {
		return zgrab2.TryGetScanStatus(err), ret, nil
	}
	if err := ret.QueryApplicationSoftwareRevision(conn); err != nil {
		return zgrab2.TryGetScanStatus(err), ret, nil
	}
	if err := ret.QueryObjectName(conn); err != nil {
		return zgrab2.TryGetScanStatus(err), ret, nil
	}
	if err := ret.QueryModelName(conn); err != nil {
		return zgrab2.TryGetScanStatus(err), ret, nil
	}
	if err := ret.QueryDescription(conn); err != nil {
		return zgrab2.TryGetScanStatus(err), ret, nil
	}
	if err := ret.QueryLocation(conn); err != nil {
		return zgrab2.TryGetScanStatus(err), ret, nil
	}

	return zgrab2.SCAN_SUCCESS, ret, nil
}
