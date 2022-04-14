package zgrab2

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zmap/zgrab2/module"
)

// All registered (available) and queued scanners (to be ran)
// are gathered and manage inside this set, similar to a module.Set.
var scans = NewScanSet()

// Scanner is an interface that represents all functions necessary to run a scan.
type Scanner interface {
	// Init runs once for this module at library init time
	Init(data module.Scan) error

	// InitPerSender runs once per Goroutine. A single Goroutine will scan some non-deterministic
	// subset of the input scan targets
	InitPerSender(senderID int) error

	// Returns the name passed at init
	GetName() string

	// Returns the trigger passed at init
	GetTrigger() string

	// Protocol returns the protocol identifier for the scan.
	Protocol() string

	// Scan connects to a host. The result should be JSON-serializable
	Scan(t ScanTarget) (ScanStatus, interface{}, error)
}

// ScanSet is the equivalent of a module set, except that it
// handles and managed both registered and running scans.
type ScanSet struct {
	// registered are the scanners bound at init time,
	// which is the entire list of scanners that can be ran,
	// including pseudo-randomly named ones generated when
	// the multiple command is used.
	registered map[string]Scanner

	// Once command is parsed, one or more scanners
	// queued is a different map used to manage scans
	// while they run. It actually has a 1-1 equivalent of
	// scansRegistered.
	queued  map[string]*Scanner
	ordered []string
}

func NewScanSet() *ScanSet {
	return &ScanSet{
		registered: make(map[string]Scanner),
		queued:     make(map[string]*Scanner),
	}
}

// Register is the simplest function to use for registering a scanner
// and its configuration flags (which are then used as CLI front-end).
func Register(name string, port int, mod module.Scan, scan Scanner) {
	module.AddModule(name, port, mod)
	scans.registered[name] = scan
}

// ScanResponse is the result of a scan on a single host.
type ScanResponse struct {
	// Status is required for all responses.
	Status ScanStatus `json:"status"`

	// Protocol is the identifier if the protocol that did the scan. In the case of a complex scan, this may differ from
	// the scan name.
	Protocol string `json:"protocol"`

	Result    interface{} `json:"result,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Error     *string     `json:"error,omitempty"`
}

// func GetScanner(name string) Scanner {
//         scan := scansQueued[name]
//         instance := reflect.New(reflect.TypeOf(scan))
//         return instance.Interface().(Scanner)
// }

// RegisterScan registers each individual scanner to be ran by the framework,
// igniting with their corresponding struct flag configurations.
func RegisterScan(name string, data module.Scan) {
	// Validate the framework configuration before we apply any
	// command-line flags to it. This is needed for special functions
	// used in scanners, and other stuff.
	ValidateFrameworkConfiguration()

	// We don't allow twice the same scanner name as key
	if scans.queued[name] != nil {
		log.Fatalf("name: %s already used", name)
	}

	// We find the scanner in the map
	scanner := scans.registered[name]

	// Initialize the scanner our configuration
	scanner.Init(data)

	scans.ordered = append(scans.ordered, name)
	scans.queued[name] = &scanner
}

// PrintScanners prints all registered scanners.
func PrintScanners() {
	for k, v := range scans.queued {
		fmt.Println(k, v)
	}
}

// RunScanner runs a single scan on a target and returns the resulting data.
func RunScanner(s Scanner, mon *Monitor, target ScanTarget) (string, ScanResponse) {
	t := time.Now()
	status, res, e := s.Scan(target)
	var err *string
	if e == nil {
		mon.statusesChan <- moduleStatus{name: s.GetName(), st: statusSuccess}
		err = nil
	} else {
		mon.statusesChan <- moduleStatus{name: s.GetName(), st: statusFailure}
		errString := e.Error()
		err = &errString
	}
	resp := ScanResponse{Result: res, Protocol: s.Protocol(), Error: err, Timestamp: t.Format(time.RFC3339), Status: status}
	return s.GetName(), resp
}

func removeTrailingNonce(s string) string {
	var j int
	for j = len(s) - 1; j > 0; j-- {
		if s[j] != '-' {
			break
		}
	}
	return s[:j+1]
}
