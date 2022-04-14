package module

import "time"

// Scan is an interface which must be implemented by all types sent to
// the flag parser.
type Scan interface {
	// Description returns a string suitable for use as an overview of this
	// module within usage text.
	Description() string

	// Help optionally returns any additional help text, e.g. specifying what empty defaults
	// are interpreted as.
	Help() string

	// Execute actually validates/enforces all command-line flags and positional arguments have valid values.
	Execute(args []string) error
}

// Base contains the options that every flags type must embed.
type Base struct {
	Port           uint          `short:"p" long:"port" description:"Specify port to grab on"`
	Name           string        `short:"n" long:"name" description:"Specify name for output json, only necessary if scanning multiple modules"`
	Timeout        time.Duration `short:"t" long:"timeout" description:"Set connection timeout (0 = no timeout)" default:"10s"`
	Trigger        string        `short:"g" long:"trigger" description:"Invoke only on targets with specified tag"`
	BytesReadLimit int           `short:"m" long:"maxbytes" description:"Maximum byte read limit per scan (0 = defaults)"`
}

// GetName returns the name of the respective scanner.
func (b *Base) GetName() string {
	return b.Name
}

// UDP contains the common options used for all UDP scans.
type UDP struct {
	LocalPort    uint   `long:"local-port" description:"Set an explicit local port for UDP traffic"`
	LocalAddress string `long:"local-addr" description:"Set an explicit local address for UDP traffic"`
}
