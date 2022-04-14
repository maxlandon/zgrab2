// Package banner provides simple banner grab and matching implementation of the zgrab2.Module.
// It sends a customizable probe (default to "\n") and filters the results based on custom regexp (--pattern)
package banner

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strconv"

	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/module"
)

// Scanner is the implementation of the zgrab2.Scanner interface.
type Scanner struct {
	config *module.Banner
	regex  *regexp.Regexp
	probe  []byte
}

// Results instances are returned by the module's Scan function.
type Results struct {
	Banner string `json:"banner,omitempty"`
	Length int    `json:"length,omitempty"`
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
	return "banner"
}

// InitPerSender initializes the scanner for a given sender.
func (scanner *Scanner) InitPerSender(senderID int) error {
	return nil
}

// Init initializes the Scanner with the command-line flags.
func (scanner *Scanner) Init(data module.Scan) error {
	var err error
	f, _ := data.(*module.Banner)
	scanner.config = f
	scanner.regex = regexp.MustCompile(scanner.config.Pattern)
	if len(f.ProbeFile) != 0 {
		scanner.probe, err = ioutil.ReadFile(f.ProbeFile)
		if err != nil {
			log.Fatal("Failed to open probe file")
			return zgrab2.ErrInvalidArguments
		}
	} else {
		strProbe, err := strconv.Unquote(fmt.Sprintf(`"%s"`, scanner.config.Probe))
		if err != nil {
			panic("Probe error")
		}
		scanner.probe = []byte(strProbe)
	}

	return nil
}

var NoMatchError = errors.New("pattern did not match")

func (scanner *Scanner) Scan(target zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	try := 0
	var (
		conn    net.Conn
		tlsConn *zgrab2.TLSConnection
		err     error
		readerr error
	)
	// Wrap the configuration with the appropriate utility code
	config := (*zgrab2.TLSFlags)(&scanner.config.TLSFlags)

	// And get the connections.
	for try < scanner.config.MaxTries {
		try++
		conn, err = target.Open(&scanner.config.Base)
		if err != nil {
			continue
		}
		if scanner.config.UseTLS {
			tlsConn, err = config.GetTLSConnection(conn)
			if err != nil {
				continue
			}
			if err = tlsConn.Handshake(); err != nil {
				continue
			}
			conn = tlsConn
		}

		break
	}
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	defer conn.Close()

	var ret []byte
	try = 0
	for try < scanner.config.MaxTries {
		try++
		_, err = conn.Write(scanner.probe)
		ret, readerr = zgrab2.ReadAvailable(conn)
		if err != nil {
			continue
		}
		if readerr != io.EOF && readerr != nil {
			continue
		}
		break
	}
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	if readerr != io.EOF && readerr != nil {
		return zgrab2.TryGetScanStatus(readerr), nil, readerr
	}
	var results Results
	if scanner.config.Hex {
		results = Results{Banner: hex.EncodeToString(ret), Length: len(ret)}
	} else {
		results = Results{Banner: string(ret), Length: len(ret)}
	}
	if scanner.regex.Match(ret) {
		return zgrab2.SCAN_SUCCESS, &results, nil
	}

	return zgrab2.SCAN_PROTOCOL_ERROR, &results, NoMatchError
}
