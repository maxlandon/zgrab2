package zgrab2

// Scanner is an interface that represents all functions necessary to run a scan
type Scanner interface {
	// Init runs once for this module at library init time
	Init(flags ScanFlags) error

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

// ScanResponse is the result of a scan on a single host
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

// ScanModule is an interface which represents a module that the framework can
// manipulate
type ScanModule interface {
	// NewFlags is called by the framework to pass to the argument parser. The parsed flags will be passed
	// to the scanner created by NewScanner().
	NewFlags() interface{}

	// NewScanner is called by the framework for each time an individual scan is specified in the config or on
	// the command-line. The framework will then call scanner.Init(name, flags).
	NewScanner() Scanner

	// Description returns a string suitable for use as an overview of this
	// module within usage text.
	Description() string
}

// ScanFlags is an interface which must be implemented by all types sent to
// the flag parser
type ScanFlags interface {
	// Help optionally returns any additional help text, e.g. specifying what empty defaults
	// are interpreted as.
	Help() string

	// Execute actually validates/enforces all command-line flags and positional arguments have valid values.
	Execute(args []string) error
}

// BaseFlags contains the options that every flags type must embed
// type BaseFlags flags.BaseFlags
// type BaseFlags flags.BaseFlags

// UDPFlags contains the common options used for all UDP scans
// type UDPFlags flags.UDPFlags

// GetName returns the name of the respective scanner
// func (b *BaseFlags) GetName() string {
//         return b.Name
// }

// GetModule returns the registered module that corresponds to the given name
// or nil otherwise
func GetModule(name string) ScanModule {
	return modules[name]
}

var modules map[string]ScanModule

func init() {
	modules = make(map[string]ScanModule)
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
