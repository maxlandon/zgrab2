package module

// Flags holds the command-line configuration for the HTTP scan module.
// Populated by the framework.
type Oracle struct {
	Base     `group:"base"`
	TLSFlags `group:"TLS"`

	// Version is the client version number sent to the server in the Connect
	// packet. TODO: Find version number mappings.
	Version uint16 `long:"client-version" description:"The client version number to send." default:"312"`

	// MinVersion is the minimum protocol version that the client claims support
	// for in the Connect packet. Same format as Version above.
	MinVersion uint16 `long:"min-server-version" description:"The minimum supported client version to send in the connect packet." default:"300"`

	// ReleaseVersion is the five-component dotted-decimal release version
	// string the client should send during native Native Security Negotiation.
	ReleaseVersion string `long:"release-version" description:"The dotted-decimal release version used during the NSN negoatiation. Must contain five components (e.g. 1.2.3.4.5)." default:"11.2.0.4.0"`

	// GlobalServiceOptions sets the ServiceOptions flags the client will send
	// to the server in the Connect packet. 16 bits.
	GlobalServiceOptions string `long:"global-service-options" description:"The Global Service Options flags to send in the connect packet." default:"0x0C41"`

	// SDU sets the requested Session Data Unit size value the client sends in
	// the Connect packet. 16 bits.
	SDU string `long:"sdu" description:"The SDU value to send in the connect packet." default:"0x2000"`

	// TDU sets the request Transport Data Unit size value the client sends in
	// the Connect packet. 16 bits.
	TDU string `long:"tdu" description:"The TDU value to send in the connect packet." default:"0xFFFF"`

	// ProtocolCharacteristics sets the protocol characteristics flags the
	// client sends to the server in the Connect packet. 16 bits.
	ProtocolCharacterisics string `long:"protocol-characteristics" description:"The Protocol Characteristics flags to send in the connect packet." default:"0x7F08"`

	// ConnectFlags sets the connect flags the client sends to the server in the
	// Connect packet. The upper 16 bits give the first byte, the lower 16 bits
	// the second byte.
	ConnectFlags string `long:"connect-flags" description:"The connect flags for the connect packet." default:"0x4141"`

	// ConnectDescriptor sets the connect descriptor the client sends in the
	// data payload of the Connect packet.
	// See https://docs.oracle.com/cd/E11882_01/network.112/e41945/glossary.htm#BGBEAGEA
	ConnectDescriptor string `long:"connect-descriptor" description:"The connect descriptor to use in the connect packet."`

	// TCPS determines whether the connection starts with a TLS handshake.
	TCPS bool `long:"tcps" description:"Wrap the connection with a TLS handshake."`

	// NewTNS causes the client to use the newer TNS header format with 32-bit
	// lengths.
	NewTNS bool `long:"new-tns" description:"If set, use new-style TNS headers"`

	// Verbose causes more verbose logging, and includes debug fields inthe scan
	// results.
	Verbose bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *Oracle) Description() string {
	return "Perform a handshake with Oracle database servers"
}

// Help returns the module's help string.
func (flags *Oracle) Help() string {
	return ""
}

// Validate checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *Oracle) Execute(args []string) error {
	return nil
}
