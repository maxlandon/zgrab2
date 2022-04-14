package module

// Flags holds the command-line configuration for the modbus scan module.
// Populated by the framework.
type Modbus struct {
	Base `group:"base"`
	// Protocols that support TLS should include zgrab2.TLSFlags
	UnitID    uint8  `long:"unit-id" description:"The UnitID / Station ID to probe"`
	ObjectID  uint8  `long:"object-id" description:"The ObjectID of the object to be read." default:"0x00"`
	Strict    bool   `long:"strict" description:"If set, perform stricter checks on the response data to get fewer false positives"`
	RequestID uint16 `long:"request-id" description:"Override the default request ID." default:"0x5A47"`
	Verbose   bool   `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *Modbus) Description() string {
	return "Probe for Modbus devices, usually PLCs as part of a SCADA system"
}

// Help returns the module's help string.
func (flags *Modbus) Help() string {
	return ""
}

// Execute checks that the flags are valid.
// On success, returns nil.
// On failure, returns an error instance describing the error.
func (flags *Modbus) Execute(args []string) error {
	if flags.Verbose {
		// If --verbose is set, do some extra checking but don't fail.
		if flags.ObjectID >= 0x07 && flags.ObjectID < 0x80 {
			// log.Warnf("ObjectIDs 0x07...0x7F are reserved (requested 0x%02x)", flags.ObjectID)
		}
	}
	return nil
}
