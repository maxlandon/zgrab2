package module

// Banner give the command-line flags for the banner module.
type Banner struct {
	Base      `group:"base"`
	TLSFlags  `group:"TLS"`
	Probe     string `long:"probe" default:"\\n" description:"Probe to send to the server. Use triple slashes to escape, for example \\\\\\n is literal \\n. Mutually exclusive with --probe-file" `
	ProbeFile string `long:"probe-file" description:"Read probe from file as byte array (hex). Mutually exclusive with --probe"`
	Pattern   string `long:"pattern" description:"Pattern to match, must be valid regexp."`
	UseTLS    bool   `long:"tls" description:"Sends probe with TLS connection. Loads TLS module command options. "`
	MaxTries  int    `long:"max-tries" default:"1" description:"Number of tries for timeouts and connection errors before giving up. Includes making TLS connection if enabled."`
	Hex       bool   `long:"hex" description:"Store banner value in hex. "`
}

// Description returns an overview of this module.
func (f *Banner) Description() string {
	return "Fetch a raw banner by sending a static probe and checking the result against a regular expression"
}

// Help returns the module's help string.
func (f *Banner) Help() string {
	return ""
}

// Execute - Actually validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (f *Banner) Execute(args []string) error {
	if f.Probe != "\\n" && f.ProbeFile != "" {
		// log.Fatal("Cannot set both --probe and --probe-file")
		// return zgrab2.ErrInvalidArguments
	}
	return nil
}
