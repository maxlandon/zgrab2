package module

// Flags holds the command-line configuration for the Telnet scan module.
// Populated by the framework.
type Telnet struct {
	Base        `group:"base"`
	MaxReadSize int  `long:"max-read-size" description:"Set the maximum number of bytes to read when grabbing the banner" default:"65536"`
	Banner      bool `long:"force-banner" description:"Always return banner if it has non-zero bytes"`
	Verbose     bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *Telnet) Description() string {
	return "Fetch a telnet banner"
}

// Help returns the module's help string.
func (flags *Telnet) Help() string {
	return ""
}

// Validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *Telnet) Execute(args []string) error {
	return nil
}
