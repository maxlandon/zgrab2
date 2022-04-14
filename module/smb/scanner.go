// Package smb provides a zgrab2 module that scans for smb.
// This was ported directly from zgrab.
package smb

import (
	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/lib/smb/smb"
	"github.com/zmap/zgrab2/module"
)

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *module.SMB
}

// Init initializes the Scanner.
func (scanner *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.SMB)
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
	return "smb"
}

// Scan performs the following:
// 1. Connect to the TCP port (default 445).
// 2. Send a negotiation packet with the default values:
//      Dialects = { DialectSmb_2_1 },
//      SecurityMode = SecurityModeSigningEnabled
// 3. Read response from server; on failure, exit with log = nil.
//      If the server returns a protocol ID indicating support for version 1, set smbv1_support = true
//      Pull out the relevant information from the response packet
// 4. If --setup-session is not set, exit with success.
// 5. Send a setup session packet to the server with appropriate values
// 6. Read the response from the server; on failure, exit with the log so far.
// 7. Return the log.
func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	conn, err := target.Open(&scanner.config.Base)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer conn.Close()
	var result *smb.SMBLog
	setupSession := scanner.config.SetupSession
	verbose := scanner.config.Verbose
	result, err = smb.GetSMBLog(conn, setupSession, false, verbose)
	if err != nil {
		if result == nil {
			conn.Close()
			conn, err = target.Open(&scanner.config.Base)
			if err != nil {
				return zgrab2.TryGetScanStatus(err), nil, err
			}
			defer conn.Close()
			result, err = smb.GetSMBLog(conn, setupSession, true, verbose)
			if err != nil {
				return zgrab2.TryGetScanStatus(err), result, err
			}
		} else {
			return zgrab2.TryGetScanStatus(err), result, err
		}
	}
	return zgrab2.SCAN_SUCCESS, result, nil
}
