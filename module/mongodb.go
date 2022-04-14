package module

// func init() {
//         mongodb.RegisterModule()
// }

// Flags contains mongodb-specific command-line flags.
type MongoDB struct {
	Base `group:"base"`
}

// Description returns an overview of this module.
func (flags *MongoDB) Description() string {
	return "Perform a handshake with a MongoDB server"
}

// Help returns the module's help string.
func (flags *MongoDB) Help() string {
	return ""
}

// Execute checks that the flags are valid.
func (flags *MongoDB) Execute(args []string) error {
	return nil
}
