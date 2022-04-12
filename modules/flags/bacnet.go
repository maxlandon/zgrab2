package flags

// Bacnet holds the command-line configuration for the bacnet scan module.
// Populated by the framework.
type Bacnet struct {
	BaseFlags `group:"base"`
	UDPFlags  `group:"udp"`
	Verbose   bool `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}
