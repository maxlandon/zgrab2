package module

// func init() {
//         mysql.RegisterModule()
// }

// Flags give the command-line flags for the MySQL module.
type MySQL struct {
	Base     `group:"base"`
	TLSFlags `group:"TLS"`
	Verbose  bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (f *MySQL) Description() string {
	return "Perform a handshake with a MySQL database"
}

// Help returns the module's help string.
func (f *MySQL) Help() string {
	return ""
}

// Execute validates the flags and returns nil on success.
func (f *MySQL) Execute(args []string) error {
	return nil
}
