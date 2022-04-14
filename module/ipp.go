package module

// Flags holds the command-line configuration for the ipp scan module.
// Populated by the framework.
type IPP struct {
	Base     `group:"base"`
	TLSFlags `group:"TLS"`
	Verbose  bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`

	// FIXME: Borrowed from http module, determine whether this is all needed
	MaxSize      int    `long:"max-size" default:"256" description:"Max kilobytes to read in response to an IPP request"`
	MaxRedirects int    `long:"max-redirects" default:"0" description:"Max number of redirects to follow"`
	UserAgent    string `long:"user-agent" default:"Mozilla/5.0 zgrab/0.x" description:"Set a custom user agent"`
	TLSRetry     bool   `long:"ipps-retry" description:"If the initial request using TLS fails, reconnect and try using plaintext IPP."`

	// FollowLocalhostRedirects overrides the default behavior to return
	// ErrRedirLocalhost whenever a redirect points to localhost.
	FollowLocalhostRedirects bool `long:"follow-localhost-redirects" description:"Follow HTTP redirects to localhost"`

	// TODO: Maybe separately implement both an ipps connection and upgrade to https
	IPPSecure bool `long:"ipps" description:"Perform a TLS handshake immediately upon connecting."`
}

// Description returns an overview of this module.
func (flags *IPP) Description() string {
	return "Probe for printers via IPP"
}

// Help returns the module's help string.
func (flags *IPP) Help() string {
	// TODO: Write a help string
	return ""
}

// Execute checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *IPP) Execute(args []string) error {
	return nil
}
