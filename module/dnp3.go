package module

// func init() {
//         dnp3.RegisterModule()
// }

// Flags holds the command-line configuration for the dnp3 scan module.
// Populated by the framework.
type DNP3 struct {
	Base `group:"base"`
	// TODO: Support UDP?
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *DNP3) Description() string {
	return "Probe for DNP3, a SCADA protocol"
}

// Help returns the module's help string.
func (flags *DNP3) Help() string {
	return ""
}

// Validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *DNP3) Execute(args []string) error {
	return nil
}
