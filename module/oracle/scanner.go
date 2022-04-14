// Package oracle provides the zgrab2 scanner module for Oracle's TNS protocol.
// Default Port: 1521 (TCP)
//
// The scan does the first part of a TNS handshake, prior to the point where
// any actual authentication is required; the happy case goes
// 1. client-to-server: Connect(--client-version, --min-server-version, --connect-descriptor)
// 2. server-to-client: Resend
// 3. client-to-server: Connect(exact same data)
// 4. server-to-client: Accept(server_version)
// 5. client-to-server: Data: Native Service Negotiation
// 6. server-to-client: Data: Native Service Negotiation(component release versions)
//
// The default scan uses a generic connect descriptor with no explicit connect
// data / service name, so it relies on the server to choose the destination.
//
// Sending an intentionally invalid --connect-descriptor can force a Refuse
// response, which should include a version number.
//
// The output includes the server's protocol version and any component release
// versions that are returned.
package oracle

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// ScanResults instances are returned by the module's Scan function.
type ScanResults struct {
	// Handshake is the log of the TNS handshake between client and server.
	Handshake *HandshakeLog `json:"handshake,omitempty"`

	// TLSLog contains the log of the TLS handshake (and any additional
	// configured TLS scan operations).
	TLSLog *zgrab2.TLSLog `json:"tls,omitempty"`
}

// Flags is a wrapper around our base declaration of Oracle flags.
// We reimplement the module.Scan interface with a new Execute function,
// so we can keep all Oracle-specific code in this package.
type Flags module.Oracle

// Description returns an overview of this module.
func (flags *Flags) Description() string {
	return "Perform a handshake with Oracle database servers"
}

// Help returns the module's help string.
func (flags *Flags) Help() string {
	return ""
}

// Validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *Flags) Execute(args []string) error {
	u16Strings := map[string]string{
		"global-service-options":   flags.GlobalServiceOptions,
		"protocol-characteristics": flags.ProtocolCharacterisics,
		"connect-flags":            flags.ConnectFlags,
		"sdu":                      flags.SDU,
		"tdu":                      flags.TDU,
	}
	for name, value := range u16Strings {
		v, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return fmt.Errorf("%s: %s is not a valid 16-bit integer: %v", name, value, err)
		}
		if v > 0xffff {
			return fmt.Errorf("%s: %s is larger than 16 bits", name, value)
		}
	}
	if _, err := EncodeReleaseVersion(flags.ReleaseVersion); err != nil {
		return fmt.Errorf("release-version: %s is not a valid five-component dotted-decimal number", flags.ReleaseVersion)
	}
	return nil
}

// Scanner implements the zgrab2.Scanner interface.
type Scanner struct {
	config *Flags
}

// Init initializes the Scanner.
func (scanner *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.Oracle)
	scanner.config = (*Flags)(f)
	if f.Verbose {
		log.SetLevel(log.DebugLevel)
	}
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
	return "oracle"
}

func (scanner *Scanner) getTNSDriver() *TNSDriver {
	mode := TNSModeOld
	if scanner.config.NewTNS {
		mode = TNSMode12c
	}
	return &TNSDriver{Mode: mode}
}

// Scan does the following:
//  1. Make a TCP connection to the target
//  2. If --tcps is set, do a TLS handshake and use the wrapped socket in future
//     calls.
//  3. Instantiate the TNS driver (TNSMode12c if --new-tns is set, otherwise
//     TNSModeOld)
//  4. Send the Connect packet to the server with the provided options and
//     connect descriptor
//  5. If the server responds with a valid TNS packet, an Oracle server has been
//     detected. If not, fail.
//  6. If the response is...
//     a. ...a Resend packet, then set result.DidResend and re-send the packet.
//     b. ...a Refused packet, then set the result.RefuseReason and RefuseError,
//        then exit.
//     c. ...a Redirect packet, then set result.RedirectTarget and exit.
//     d. ...an Accept packet, go to 7
//     e. ...anything else: exit with SCAN_APPLICATION_ERROR
//  7. Pull the server protocol version and other flags from the Accept packet
//     into the results, then send a Native Security Negotiation Data packet.
//  8. If the response is not a Data packet, exit with SCAN_APPLICATION_ERROR.
//  9. Pull the versions out of the response and exit with SCAN_SUCCESS.
func (scanner *Scanner) Scan(t zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	var results *ScanResults

	sock, err := t.Open(&scanner.config.Base)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	if scanner.config.TCPS {
		// Wrap the configuration with the appropriate utility code
		config := (*zgrab2.TLSFlags)(&scanner.config.TLSFlags)

		tlsConn, err := config.GetTLSConnection(sock)
		if err != nil {
			// GetTLSConnection can only fail if the input flags are bad
			panic(err)
		}
		results = new(ScanResults)
		results.TLSLog = tlsConn.GetLog()
		err = tlsConn.Handshake()
		if err != nil {
			return zgrab2.TryGetScanStatus(err), nil, err
		}
		sock = tlsConn
	}

	conn := Connection{
		conn:      sock,
		scanner:   scanner,
		target:    &t,
		tnsDriver: scanner.getTNSDriver(),
	}
	connectDescriptor := scanner.config.ConnectDescriptor
	if connectDescriptor == "" {
		// In local testing, omitting the SERVICE_NAME allowed the server to
		// choose an appropriate default. CID.PROGRAM added strictly for logging
		// purposes.
		connectDescriptor = "(DESCRIPTION=(CONNECT_DATA=(CID=(PROGRAM=zgrab2))))"
	}
	handshakeLog, err := conn.Connect(connectDescriptor)
	if handshakeLog != nil {
		// Ensure that any handshake logs, even if incomplete, get returned.
		if results == nil {
			// If the results were not created previously to store the TLS log,
			// create it now
			results = new(ScanResults)
		}
		results.Handshake = handshakeLog
	}

	if err != nil {
		switch err {
		case ErrUnexpectedResponse:
			return zgrab2.SCAN_APPLICATION_ERROR, results, err
		default:
			return zgrab2.TryGetScanStatus(err), results, err
		}
	}

	return zgrab2.SCAN_SUCCESS, results, nil
}
