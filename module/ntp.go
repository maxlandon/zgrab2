package module

// func init() {
//         ntp.RegisterModule()
// }

// Flags holds the command-line flags for the scanner.
type NTP struct {
	Base          `group:"base"`
	UDP           `group:"UDP"`
	Verbose       bool   `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
	Version       uint8  `long:"version" description:"The version number to pass to the Server." default:"3"`
	LeapIndicator uint8  `long:"leap-indicator" description:"The LI value to pass to the Server. Default 3 (Unknown)"`
	SkipGetTime   bool   `long:"skip-get-time" description:"If set, don't request the Server time"`
	MonList       bool   `long:"monlist" description:"Perform a ReqMonGetList request"`
	RequestCode   string `long:"request-code" description:"Specify a request code for MonList other than ReqMonGetList" default:"REQ_MON_GETLIST"`
}

// Description returns an overview of this module.
func (cfg *NTP) Description() string {
	return "Scan for NTP"
}

// Validate checks that the flags are valid.
func (cfg *NTP) Execute(args []string) error {
	return nil
}

// Help returns the module's help string.
func (cfg *NTP) Help() string {
	return ""
}
