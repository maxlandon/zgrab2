package module

// Flags holds the command-line configuration for the fox scan module.
// Populated by the framework.
type Fox struct {
	Base `group:"base"`

	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *Fox) Description() string {
	return "Probe for Tridium Fox"
}

// Help returns the module's help string.
func (flags *Fox) Help() string {
	return ""
}

// Validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *Fox) Execute(args []string) error {
	return nil
}
