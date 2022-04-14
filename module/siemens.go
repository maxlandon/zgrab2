package module

// func init() {
//         siemens.RegisterModule()
// }

// Flags holds the command-line configuration for the siemens scan module.
// Populated by the framework.
type Siemens struct {
	Base `group:"base"`
	// TODO: configurable TSAP source / destination, etc
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *Siemens) Description() string {
	return "Probe for Siemens S7 devices"
}

// Validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *Siemens) Execute(args []string) error {
	return nil
}

// Help returns the module's help string.
func (flags *Siemens) Help() string {
	return ""
}
