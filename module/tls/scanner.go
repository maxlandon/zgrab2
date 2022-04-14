package tls

import (
	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

type Scanner struct {
	config *module.TLS
}

func (s *Scanner) Init(flags module.Scan) error {
	f, ok := flags.(*module.TLS)
	if !ok {
		return zgrab2.ErrMismatchedFlags
	}
	s.config = f
	return nil
}

func (s *Scanner) GetName() string {
	return s.config.Name
}

func (s *Scanner) GetTrigger() string {
	return s.config.Trigger
}

func (s *Scanner) InitPerSender(senderID int) error {
	return nil
}

// Scan opens a TCP connection to the target (default port 443), then performs
// a TLS handshake. If the handshake gets past the ServerHello stage, the
// handshake log is returned (along with any other TLS-related logs, such as
// heartbleed, if enabled).
func (s *Scanner) Scan(t zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	// We need to cast the flags first, which
	// has the effect of wrapping them with some connection utility code
	tlsFlags := (*zgrab2.TLSFlags)(&s.config.TLSFlags)
	conn, err := t.OpenTLS(&s.config.Base, tlsFlags)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		if conn != nil {
			if log := conn.GetLog(); log != nil {
				if log.HandshakeLog.ServerHello != nil {
					// If we got far enough to get a valid ServerHello, then
					// consider it to be a positive TLS detection.
					return zgrab2.TryGetScanStatus(err), log, err
				}
				// Otherwise, detection failed.
			}
		}
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	return zgrab2.SCAN_SUCCESS, conn.GetLog(), nil
}

// Protocol returns the protocol identifier for the scanner.
func (s *Scanner) Protocol() string {
	return "tls"
}
