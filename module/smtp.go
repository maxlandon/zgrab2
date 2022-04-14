package module

// import "github.com/zmap/zgrab2/modules/smtp"

// func init() {
//         smtp.RegisterModule()
// }

// Flags holds the command-line configuration for the HTTP scan module.
// Populated by the framework.
type SMTP struct {
	Base     `group:"base"`
	TLSFlags `group:"TLS"`

	// SendHELO indicates that the EHLO command should be set.
	SendEHLO bool `long:"send-ehlo" description:"Send the EHLO command; use --ehlo-domain to set a domain."`

	// SendEHLO indicates that the EHLO command should be set.
	SendHELO bool `long:"send-helo" description:"Send the EHLO command; use --helo-domain to set a domain."`

	// SendHELP indicates that the client should send the HELP command (after HELO/EHLO).
	SendHELP bool `long:"send-help" description:"Send the HELP command"`

	// SendQUIT indicates that the QUIT command should be set.
	SendQUIT bool `long:"send-quit" description:"Send the QUIT command before closing."`

	// HELODomain is the domain the client should send in the HELO command.
	HELODomain string `long:"helo-domain" description:"Set the domain to use with the HELO command. Implies --send-helo."`

	// EHLODomain is the domain the client should send in the HELO command.
	EHLODomain string `long:"ehlo-domain" description:"Set the domain to use with the EHLO command. Implies --send-ehlo."`

	// SMTPSecure indicates that the entire transaction should be wrapped in a TLS session.
	SMTPSecure bool `long:"smtps" description:"Perform a TLS handshake immediately upon connecting."`

	// StartTLS indicates that the client should attempt to update the connection to TLS.
	StartTLS bool `long:"starttls" description:"Send STARTTLS before negotiating"`

	// Verbose indicates that there should be more verbose logging.
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *SMTP) Description() string {
	return "Fetch an SMTP server banner, optionally over TLS"
}

// Help returns the module's help string.
func (flags *SMTP) Help() string {
	return ""
}

// Validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *SMTP) Execute(args []string) error {
	if flags.StartTLS && flags.SMTPSecure {
		// log.Errorln("Cannot specify both --smtps and --starttls")
		// return zgrab2.ErrInvalidArguments
	}
	if flags.EHLODomain != "" {
		flags.SendEHLO = true
	}
	if flags.HELODomain != "" {
		flags.SendHELO = true
	}
	if flags.SendHELO && flags.SendEHLO {
		// log.Errorln("Cannot provide both EHLO and HELO")
		// return zgrab2.ErrInvalidArguments
	}
	return nil
}
