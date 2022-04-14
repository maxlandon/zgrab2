package module

// func init() {
//         postgres.RegisterModule()
// }

// Flags sets the module-specific flags that can be passed in from the
// command line.
type Postgres struct {
	Base            `group:"base"`
	TLSFlags        `group:"TLS"`
	SkipSSL         bool   `long:"skip-ssl" description:"If set, do not attempt to negotiate an SSL connection"`
	Verbose         bool   `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
	ProtocolVersion string `long:"protocol-version" description:"The protocol to use in the StartupPacket" default:"3.0"`
	User            string `long:"user" description:"Username to pass to StartupMessage. If omitted, no user will be sent." default:""`
	Database        string `long:"database" description:"Database to pass to StartupMessage. If omitted, none will be sent." default:""`
	ApplicationName string `long:"application-name" description:"application_name value to pass in StartupMessage. If omitted, none will be sent." default:""`
}

// Help returns the module's help string.
func (f *Postgres) Help() string {
	return ""
}

// Description returns an overview of this module.
func (f *Postgres) Description() string {
	return "Perform a handshake with a PostgreSQL server"
}

// Validate checks the arguments; on success, returns nil.
func (f *Postgres) Execute(args []string) error {
	return nil
}
