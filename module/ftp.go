package module

import "fmt"

// Flags are the FTP-specific command-line flags. Taken from the original zgrab.
// (TODO: should FTPAuthTLS be on by default?).
type FTP struct {
	Base        `group:"base"`
	TLSFlags    `group:"TLS"`
	Verbose     bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
	FTPAuthTLS  bool `long:"authtls" description:"Collect FTPS certificates in addition to FTP banners"`
	ImplicitTLS bool `long:"implicit-tls" description:"Attempt to connect via a TLS wrapped connection"`
}

// Description returns an overview of this module.
func (m *FTP) Description() string {
	return "Grab an FTP banner"
}

// Help returns this module's help string.
func (f *FTP) Help() string {
	return ""
}

// Execute flags.
func (f *FTP) Execute(args []string) (err error) {
	if f.FTPAuthTLS && f.ImplicitTLS {
		err = fmt.Errorf("Cannot specify both '--authtls' and '--implicit-tls' together")
	}
	return
}
