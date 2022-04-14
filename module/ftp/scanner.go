// Package ftp contains the zgrab2 Module implementation for FTP(S).
//
// Setting the --authtls flag will cause the scanner to attempt a upgrade the
// connection to TLS. Settings for the TLS handshake / probe can be set with
// the standard TLSFlags.
//
// The scan performs a banner grab and (optionally) a TLS handshake.
//
// The output is the banner, any responses to the AUTH TLS/AUTH SSL commands,
// and any TLS logs.
package ftp

import (
	"net"
	"regexp"
	"strings"

	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// ScanResults is the output of the scan.
// Identical to the original from zgrab, with the addition of TLSLog.
type ScanResults struct {
	// Banner is the initial data banner sent by the server.
	Banner string `json:"banner,omitempty"`

	// AuthTLSResp is the response to the AUTH TLS command.
	// Only present if the FTPAuthTLS flag is set.
	AuthTLSResp string `json:"auth_tls,omitempty"`

	// AuthSSLResp is the response to the AUTH SSL command.
	// Only present if the FTPAuthTLS flag is set and AUTH TLS failed.
	AuthSSLResp string `json:"auth_ssl,omitempty"`

	// ImplicitTLS is true if the connection is wrapped in TLS, as opposed
	// to via AUTH TLS or AUTH SSL.
	ImplicitTLS bool `json:"implicit_tls,omitempty"`

	// TLSLog is the standard shared TLS handshake log.
	// Only present if the FTPAuthTLS flag is set.
	TLSLog *zgrab2.TLSLog `json:"tls,omitempty"`
}

// Scanner implements the zgrab2.Scanner interface, and holds the state
// for a single scan.
type Scanner struct {
	config *module.FTP
}

// Connection holds the state for a single connection to the FTP server.
type Connection struct {
	// buffer is a temporary buffer for sending commands -- so, never interleave
	// sendCommand calls on a given connection
	buffer  [10000]byte
	config  *module.FTP
	results ScanResults
	conn    net.Conn
}

// Protocol returns the protocol identifier for the scanner.
func (s *Scanner) Protocol() string {
	return "ftp"
}

// Init initializes the Scanner instance with the flags from the command
// line.
func (s *Scanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.FTP)
	s.config = f
	return nil
}

// InitPerSender does nothing in this module.
func (s *Scanner) InitPerSender(senderID int) error {
	return nil
}

// GetName returns the configured name for the Scanner.
func (s *Scanner) GetName() string {
	return s.config.Name
}

// GetTrigger returns the Trigger defined in the Flags.
func (scanner *Scanner) GetTrigger() string {
	return scanner.config.Trigger
}

// ftpEndRegex matches zero or more lines followed by a numeric FTP status code
// and linebreak, e.g. "200 OK\r\n".
var ftpEndRegex = regexp.MustCompile(`^(?:.*\r?\n)*([0-9]{3})( [^\r\n]*)?\r?\n$`)

// isOKResponse returns true iff and only if the given response code indicates
// success (e.g. 2XX).
func (ftp *Connection) isOKResponse(retCode string) bool {
	// TODO: This is the current behavior; should it check that it isn't
	// garbage that happens to start with 2 (e.g. it's only ASCII chars, the
	// prefix is 2[0-9]+, etc)?
	return strings.HasPrefix(retCode, "2")
}

// readResponse reads an FTP response chunk from the server.
// It returns the full response, as well as the status code alone.
func (ftp *Connection) readResponse() (string, string, error) {
	respLen, err := zgrab2.ReadUntilRegex(ftp.conn, ftp.buffer[:], ftpEndRegex)
	if err != nil {
		return "", "", err
	}
	ret := string(ftp.buffer[0:respLen])
	retCode := ftpEndRegex.FindStringSubmatch(ret)[1]
	return ret, retCode, nil
}

// GetFTPBanner reads the data sent by the server immediately after connecting.
// Returns true if and only if the server returns a success status code.
// Taken over from the original zgrab.
func (ftp *Connection) GetFTPBanner() (bool, error) {
	banner, retCode, err := ftp.readResponse()
	if err != nil {
		return false, err
	}
	ftp.results.Banner = banner
	return ftp.isOKResponse(retCode), nil
}

// sendCommand sends a command and waits for / reads / returns the response.
func (ftp *Connection) sendCommand(cmd string) (string, string, error) {
	ftp.conn.Write([]byte(cmd + "\r\n"))
	return ftp.readResponse()
}

// SetupFTPS returns true if and only if the server reported support for FTPS.
// First attempt AUTH TLS; if that fails, try AUTH SSL.
// Taken over from the original zgrab.
func (ftp *Connection) SetupFTPS() (bool, error) {
	ret, retCode, err := ftp.sendCommand("AUTH TLS")
	if err != nil {
		return false, err
	}
	ftp.results.AuthTLSResp = ret
	if ftp.isOKResponse(retCode) {
		return true, nil
	}
	ret, retCode, err = ftp.sendCommand("AUTH SSL")
	if err != nil {
		return false, err
	}
	ftp.results.AuthSSLResp = ret

	if ftp.isOKResponse(retCode) {
		return true, nil
	}
	return false, nil
}

// GetFTPSCertificates attempts to perform a TLS handshake with the server so
// that the TLS certificates will end up in the TLSLog.
// First sends the AUTH TLS/AUTH SSL command to tell the server we want to
// do a TLS handshake. If that fails, break. Otherwise, perform the handshake.
// Taken over from the original zgrab.
func (ftp *Connection) GetFTPSCertificates() error {
	ftpsReady, err := ftp.SetupFTPS()
	if err != nil {
		return err
	}
	if !ftpsReady {
		return nil
	}
	// Wrap the configuration with the appropriate utility code
	config := (*zgrab2.TLSFlags)(&ftp.config.TLSFlags)

	var conn *zgrab2.TLSConnection
	if conn, err = config.GetTLSConnection(ftp.conn); err != nil {
		return err
	}
	ftp.results.TLSLog = conn.GetLog()

	if err = conn.Handshake(); err != nil {
		// NOTE: With the default config of vsftp (without ssl_ciphers=HIGH),
		// AUTH TLS succeeds, but the handshake fails, dumping
		// "error:1408A0C1:SSL routines:ssl3_get_client_hello:no shared cipher"
		// to the socket.
		return err
	}
	ftp.conn = conn
	return nil
}

// Scan performs the configured scan on the FTP server, as follows:
// * Read the banner into results.Banner (if it is not a 2XX response, bail)
// * If the FTPAuthTLS flag is not set, finish.
// * Send the AUTH TLS command to the server. If the response is not 2XX, then
//   send the AUTH SSL command. If the response is not 2XX, then finish.
// * Perform ths TLS handshake / any configured TLS scans, populating
//   results.TLSLog.
// * Return SCAN_SUCCESS, &results, nil.
func (s *Scanner) Scan(t zgrab2.ScanTarget) (status zgrab2.ScanStatus, result interface{}, thrown error) {
	var err error
	conn, err := t.Open(&s.config.Base)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	cn := conn
	defer func() {
		cn.Close()
	}()

	results := ScanResults{}
	if s.config.ImplicitTLS {
		// Wrap the configuration with the appropriate utility code
		config := (*zgrab2.TLSFlags)(&s.config.TLSFlags)

		tlsConn, err := config.GetTLSConnection(conn)
		if err != nil {
			return zgrab2.TryGetScanStatus(err), nil, err
		}
		results.ImplicitTLS = true
		results.TLSLog = tlsConn.GetLog()
		err = tlsConn.Handshake()
		if err != nil {
			return zgrab2.TryGetScanStatus(err), nil, err
		}
		cn = tlsConn
	}

	ftp := Connection{conn: cn, config: s.config, results: results}
	is200Banner, err := ftp.GetFTPBanner()
	if err != nil {
		return zgrab2.TryGetScanStatus(err), &ftp.results, err
	}
	if s.config.FTPAuthTLS && is200Banner {
		if err := ftp.GetFTPSCertificates(); err != nil {
			return zgrab2.SCAN_APPLICATION_ERROR, &ftp.results, err
		}
	}
	return zgrab2.SCAN_SUCCESS, &ftp.results, nil
}
