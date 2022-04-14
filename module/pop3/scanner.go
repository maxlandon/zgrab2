// Package pop3 provides a zgrab2 module that scans for POP3 mail
// servers.
// Default Port: 110 (TCP)
//
// The --send-help and --send-noop flags tell the scanner to send a
// HELP or NOOP command and read the response.
//
// The --pop3s flag tells the scanner to perform a TLS handshake
// immediately after connecting, before even attempting to read
// the banner.
// The --starttls flag tells the scanner to send the STLS command,
// and then negotiate a TLS connection.
// The scanner uses the standard TLS flags for the handshake.
// --pop3s and --starttls are mutually exclusive.
// --pop3s does not change the default port number from 110, so
// it should usually be coupled with e.g. --port 995.
//
// The --send-quit flag tells the scanner to send a QUIT command
// before disconnecting.
//
// So, if no flags are specified, the scanner simply reads the banner
// returned by the server and disconnects.
//
// The output contains the banner and the responses to any commands that
// were sent, and if or --pop3s --starttls were set, the standard TLS logs.
package pop3

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// ScanResults instances are returned by the module's Scan function.
type ScanResults struct {
	// Banner is the string sent by the server immediately after connecting.
	Banner string `json:"banner,omitempty"`

	// NOOP is the server's response to the NOOP command, if one is sent.
	NOOP string `json:"noop,omitempty"`

	// HELP is the server's response to the HELP command, if it is sent.
	HELP string `json:"help,omitempty"`

	// StartTLS is the server's response to the STARTTLS command, if it is sent.
	StartTLS string `json:"starttls,omitempty"`

	// QUIT is the server's response to the QUIT command, if it is sent.
	QUIT string `json:"quit,omitempty"`

	// TLSLog is the standard TLS log, if --starttls or --pop3s is enabled.
	TLSLog *zgrab2.TLSLog `json:"tls,omitempty"`
}

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *module.POP3
}

// Init initializes the Scanner.
func (scanner *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.POP3)
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
	return "pop3"
}

func getPOP3Error(response string) error {
	if !strings.HasPrefix(response, "-") {
		return nil
	}
	return fmt.Errorf("POP3 error: %s", response[1:])
}

// Check the contents of the POP3 header and return a relevant ScanStatus.
func VerifyPOP3Contents(banner string) zgrab2.ScanStatus {
	lowerBanner := strings.ToLower(banner)
	switch {
	case strings.HasPrefix(banner, "-ERR "):
		return zgrab2.SCAN_APPLICATION_ERROR
	case strings.HasPrefix(banner, "+OK "),
		strings.Contains(banner, "POP3"),
		// These are rare for POP3 if they happen at all,
		// But it won't hurt to check just in case as a backup
		strings.Contains(lowerBanner, "blacklist"),
		strings.Contains(lowerBanner, "abuse"),
		strings.Contains(lowerBanner, "rbl"),
		strings.Contains(lowerBanner, "spamhaus"),
		strings.Contains(lowerBanner, "relay"):
		return zgrab2.SCAN_SUCCESS
	default:
		return zgrab2.SCAN_PROTOCOL_ERROR
	}
}

// Scan performs the POP3 scan.
// 1. Open a TCP connection to the target port (default 110).
// 2. If --pop3s is set, perform a TLS handshake using the command-line
//    flags.
// 3. Read the banner.
// 4. If --send-help is sent, send HELP, read the result.
// 5. If --send-noop is sent, send NOOP, read the result.
// 6. If --starttls is sent, send STLS, read the result, negotiate a
//    TLS connection using the command-line flags.
// 7. If --send-quit is sent, send QUIT and read the result.
// 8. Close the connection.
func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	c, err := target.Open(&scanner.config.Base)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer c.Close()

	// Wrap the configuration with the appropriate utility code
	config := (*zgrab2.TLSFlags)(&scanner.config.TLSFlags)

	result := &ScanResults{}
	if scanner.config.POP3Secure {
		tlsConn, err := config.GetTLSConnection(c)
		if err != nil {
			return zgrab2.TryGetScanStatus(err), nil, err
		}
		result.TLSLog = tlsConn.GetLog()
		if err := tlsConn.Handshake(); err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		c = tlsConn
	}
	conn := Connection{Conn: c}
	banner, err := conn.ReadResponse()
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	// Quit early if no valid response
	// OR save it to return later
	sr := VerifyPOP3Contents(banner)
	if sr == zgrab2.SCAN_PROTOCOL_ERROR {
		return sr, nil, errors.New("Invalid response for POP3")
	}
	result.Banner = banner
	if scanner.config.SendHELP {
		ret, err := conn.SendCommand("HELP")
		if err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		result.HELP = ret
	}
	if scanner.config.SendNOOP {
		ret, err := conn.SendCommand("NOOP")
		if err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		result.NOOP = ret
	}
	if scanner.config.StartTLS {
		ret, err := conn.SendCommand("STLS")
		if err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		result.StartTLS = ret
		if err := getPOP3Error(ret); err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		tlsConn, err := config.GetTLSConnection(conn.Conn)
		if err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		result.TLSLog = tlsConn.GetLog()
		if err := tlsConn.Handshake(); err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		conn.Conn = tlsConn
	}
	if scanner.config.SendQUIT {
		ret, err := conn.SendCommand("QUIT")
		if err != nil {
			if err != nil {
				return zgrab2.TryGetScanStatus(err), nil, err
			}
		}
		result.QUIT = ret
	}
	return sr, result, nil
}
