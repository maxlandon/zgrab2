// Package smtp provides a zgrab2 module that scans for SMTP mail
// servers.
// Default Port: 25 (TCP)
//
// The --smtps command tells the scanner to wrap the entire connection
// in a TLS session.
//
// The --send-ehlo and --send-helo flags tell the scanner to first send
// the EHLO/HELO command; if a --ehlo-domain or --helo-domain is present
// that domain will be used, otherwise it is omitted.
// The EHLO and HELO flags are mutually exclusive.
//
// The --send-help flag tells the scanner to send a HELP command.
//
// The --starttls flag tells the scanner to send the STARTTLS command,
// and then negotiate a TLS connection.
// The scanner uses the standard TLS flags for the handshake.
//
// The --send-quit flag tells the scanner to send a QUIT command.
//
// So, if no flags are specified, the scanner simply reads the banner
// returned by the server and disconnects.
//
// The output contains the banner and the responses to any commands that
// were sent, and if --starttls or --smtps was sent, the standard TLS logs.
package smtp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// ErrInvalidResponse is returned when the server returns an invalid or unexpected response.
var ErrInvalidResponse = zgrab2.NewScanError(zgrab2.SCAN_PROTOCOL_ERROR, errors.New("Invalid response for SMTP"))

// ScanResults instances are returned by the module's Scan function.
type ScanResults struct {
	// Banner is the string sent by the server immediately after connecting.
	Banner string `json:"banner,omitempty"`

	// HELO is the server's response to the HELO command, if one is sent.
	HELO string `json:"helo,omitempty"`

	// EHLO is the server's response to the EHLO command, if one is sent.
	EHLO string `json:"ehlo,omitempty"`

	// HELP is the server's response to the HELP command, if it is sent.
	HELP string `json:"help,omitempty"`

	// StartTLS is the server's response to the STARTTLS command, if it is sent.
	StartTLS string `json:"starttls,omitempty"`

	// QUIT is the server's response to the QUIT command, if it is sent.
	QUIT string `json:"quit,omitempty"`

	// ImplicitTLS is true if the connection was wrapped in TLS, as opposed
	// to using StartTls
	ImplicitTLS bool `json:"implicit_tls,omitempty"`

	// TLSLog is the standard TLS log, if STARTTLS is sent.
	TLSLog *zgrab2.TLSLog `json:"tls,omitempty"`
}

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *module.SMTP
}

// Init initializes the Scanner.
func (scanner *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.SMTP)
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
	return "smtp"
}

func getSMTPCode(response string) (int, error) {
	if len(response) < 5 {
		return 0, ErrInvalidResponse
	}
	ret, err := strconv.Atoi(response[0:3])
	if err != nil {
		return 0, ErrInvalidResponse
	}
	return ret, nil
}

// Get a command with an optional argument (so if the argument is absent, there is no trailing space).
func getCommand(cmd string, arg string) string {
	if arg == "" {
		return cmd
	}
	return cmd + " " + arg
}

// Verify that an SMTP code was returned, and that it is a successful one!
// Return code on SCAN_APPLICATION_ERROR for better info.
func VerifySMTPContents(banner string) (zgrab2.ScanStatus, int) {
	code, err := getSMTPCode(banner)
	lowerBanner := strings.ToLower(banner)
	switch {
	case err == nil && (code < 200 || code >= 300):
		return zgrab2.SCAN_APPLICATION_ERROR, code
	case err == nil,
		strings.Contains(banner, "SMTP"),
		strings.Contains(lowerBanner, "blacklist"),
		strings.Contains(lowerBanner, "abuse"),
		strings.Contains(lowerBanner, "rbl"),
		strings.Contains(lowerBanner, "spamhaus"),
		strings.Contains(lowerBanner, "relay"):
		return zgrab2.SCAN_SUCCESS, 0
	default:
		return zgrab2.SCAN_PROTOCOL_ERROR, 0
	}
}

// Scan performs the SMTP scan.
// 1. Open a TCP connection to the target port (default 25).
// 2. If --smtps is set, perform a TLS handshake.
// 3. Read the banner.
// 4. If --send-ehlo or --send-helo is sent, send the corresponding EHLO
//    or HELO command.
// 5. If --send-help is sent, send HELP, read the result.
// 6. If --starttls is sent, send STARTTLS, read the result, negotiate a
//    TLS connection.
// 7. If --send-quit is sent, send QUIT and read the result.
// 8. Close the connection.
func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	c, err := target.Open(&scanner.config.Base)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer c.Close()
	result := &ScanResults{}

	// Wrap the configuration with the appropriate utility code
	config := (*zgrab2.TLSFlags)(&scanner.config.TLSFlags)

	if scanner.config.SMTPSecure {
		tlsConn, err := config.GetTLSConnection(c)
		if err != nil {
			return zgrab2.TryGetScanStatus(err), nil, err
		}
		result.TLSLog = tlsConn.GetLog()
		if err := tlsConn.Handshake(); err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		c = tlsConn
		result.ImplicitTLS = true
	}
	conn := Connection{Conn: c}
	banner, err := conn.ReadResponse()
	if err != nil {
		if !scanner.config.SMTPSecure {
			result = nil
		}
		return zgrab2.TryGetScanStatus(err), result, err
	}
	// Quit early if we didn't get a valid response
	// OR save response to return later
	sr, bannerResponseCode := VerifySMTPContents(banner)
	if sr == zgrab2.SCAN_PROTOCOL_ERROR {
		return sr, nil, errors.New("Invalid response for SMTP")
	}
	result.Banner = banner
	if scanner.config.SendHELO {
		ret, err := conn.SendCommand(getCommand("HELO", scanner.config.HELODomain))
		if err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		result.HELO = ret
	}
	if scanner.config.SendEHLO {
		ret, err := conn.SendCommand(getCommand("EHLO", scanner.config.EHLODomain))
		if err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		result.EHLO = ret
	}
	if scanner.config.SendHELP {
		ret, err := conn.SendCommand("HELP")
		if err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		result.HELP = ret
	}
	if scanner.config.StartTLS {
		ret, err := conn.SendCommand("STARTTLS")
		if err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		result.StartTLS = ret
		code, err := getSMTPCode(ret)
		if err != nil {
			return zgrab2.TryGetScanStatus(err), result, err
		}
		if code < 200 || code >= 300 {
			return zgrab2.SCAN_APPLICATION_ERROR, result, fmt.Errorf("SMTP error code %d returned from STARTTLS command (%s)", code, strings.TrimSpace(ret))
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
	if sr == zgrab2.SCAN_APPLICATION_ERROR {
		return sr, result, fmt.Errorf("SMTP error code %d returned in banner grab", bannerResponseCode)
	}
	return sr, result, nil
}
