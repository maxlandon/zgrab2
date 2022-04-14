package module

// Bacnet holds the command-line configuration for the bacnet scan module.
// Populated by the framework.
type Bacnet struct {
	Base    `group:"base"`
	UDP     `group:"udp"`
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Execute - Actually validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (f *Bacnet) Execute(args []string) error {
	return nil
}

// Help returns the module's help string.
func (f *Bacnet) Help() string {
	return ""
}

// Description returns text uses in the help for this module.
func (f *Bacnet) Description() string {
	return "Probe for devices that speak Bacnet, commonly used for HVAC control."
}
