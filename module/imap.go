package module

// Flags holds the command-line configuration for the IMAP scan module.
// Populated by the framework.
type IMAP struct {
	Base     `group:"base"`
	TLSFlags `group:"TLS"`

	// SendCLOSE indicates that the CLOSE command should be sent.
	SendCLOSE bool `long:"send-close" description:"Send the CLOSE command before closing."`

	// IMAPSecure indicates that the client should do a TLS handshake immediately after connecting.
	IMAPSecure bool `long:"imaps" description:"Immediately negotiate a TLS connection"`

	// StartTLS indicates that the client should attempt to update the connection to TLS.
	StartTLS bool `long:"starttls" description:"Send STLS before negotiating"`

	// Verbose indicates that there should be more verbose logging.
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *IMAP) Description() string {
	return "Fetch an IMAP banner, optionally over TLS"
}

// Help returns the module's help string.
func (flags *IMAP) Help() string {
	return ""
}

// Execute checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *IMAP) Execute(args []string) error {
	if flags.StartTLS && flags.IMAPSecure {
		// log.Error("Cannot send both --starttls and --imaps")
		// return zgrab2.ErrInvalidArguments
	}
	return nil
}
