package module

// Flags give the command-line flags for the banner module.
type JARM struct {
	Base     `group:"base"`
	MaxTries int `long:"max-tries" default:"1" description:"Number of tries for timeouts and connection errors before giving up."`
}

// Description returns an overview of this module.
func (flags *JARM) Description() string {
	return "Send TLS requiests and generate a JARM fingerprint"
}

// Help returns the module's help string.
func (f *JARM) Help() string {
	return ""
}

// Execute validates the flags and returns nil on success.
func (f *JARM) Execute(args []string) error {
	return nil
}
