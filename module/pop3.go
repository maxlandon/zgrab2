package module

// func init() {
//         pop3.RegisterModule()
// }

// Flags holds the command-line configuration for the POP3 scan module.
// Populated by the framework.
type POP3 struct {
	Base     `group:"base"`
	TLSFlags `group:"TLS"`

	// SendHELP indicates that the client should send the HELP command.
	SendHELP bool `long:"send-help" description:"Send the HELP command"`

	// SendNOOP indicates that the NOOP command should be sent.
	SendNOOP bool `long:"send-noop" description:"Send the NOOP command before closing."`

	// SendQUIT indicates that the QUIT command should be sent.
	SendQUIT bool `long:"send-quit" description:"Send the QUIT command before closing."`

	// POP3Secure indicates that the client should do a TLS handshake immediately after connecting.
	POP3Secure bool `long:"pop3s" description:"Immediately negotiate a TLS connection"`

	// StartTLS indicates that the client should attempt to update the connection to TLS.
	StartTLS bool `long:"starttls" description:"Send STLS before negotiating"`

	// Verbose indicates that there should be more verbose logging.
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *POP3) Description() string {
	return "Fetch POP3 banners, optionally over TLS"
}

// Help returns the module's help string.
func (flags *POP3) Help() string {
	return ""
}

// Validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *POP3) Execute(args []string) error {
	if flags.StartTLS && flags.POP3Secure {
		// log.Error("Cannot send both --starttls and --pop3s")
		// return zgrab2.ErrInvalidArguments
	}
	return nil
}
