package module

// func init() {
//         mssql.RegisterModule()
// }

// Flags defines the command-line configuration options for the module.
type MSSQL struct {
	Base        `group:"base"`
	TLSFlags    `group:"TLS"`
	EncryptMode string `long:"encrypt-mode" description:"The type of encryption to request in the pre-login step. One of ENCRYPT_ON, ENCRYPT_OFF, ENCRYPT_NOT_SUP." default:"ENCRYPT_ON"`
	Verbose     bool   `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *MSSQL) Description() string {
	return "Perform a handshake for MSSQL databases"
}

// Help returns the help string for this module.
func (flags *MSSQL) Help() string {
	return ""
}

// Execute does nothing in this module.
func (flags *MSSQL) Execute(args []string) error {
	return nil
}
