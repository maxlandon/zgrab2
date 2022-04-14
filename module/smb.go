package module

// func init() {
//         smb.RegisterModule()
// }

// Flags holds the command-line configuration for the smb scan module.
// Populated by the framework.
type SMB struct {
	Base `group:"base"`

	// SetupSession tells the client to continue the handshake up to the point where credentials would be needed.
	SetupSession bool `long:"setup-session" description:"After getting the response from the negotiation request, send a setup session packet."`

	// Verbose requests more verbose logging / output.
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *SMB) Description() string {
	return "Probe for SMB servers (Windows filesharing / SAMBA)"
}

// Help returns the module's help string.
func (flags *SMB) Help() string {
	return ""
}

// Validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *SMB) Execute(args []string) error {
	return nil
}
