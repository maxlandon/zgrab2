package module

// func init() {
//         redis.RegisterModule()
// }

// Flags contains redis-specific command-line flags.
type Redis struct {
	Base `group:"base"`

	CustomCommands   string `long:"custom-commands" description:"Pathname for JSON/YAML file that contains extra commands to execute. WARNING: This is sent in the clear."`
	Mappings         string `long:"mappings" description:"Pathname for JSON/YAML file that contains mappings for command names."`
	MaxInputFileSize int64  `long:"max-input-file-size" default:"102400" description:"Maximum size for either input file."`
	Password         string `long:"password" description:"Set a password to use to authenticate to the server. WARNING: This is sent in the clear."`
	DoInline         bool   `long:"inline" description:"Send commands using the inline syntax"`
	Verbose          bool   `long:"verbose" description:"More verbose logging, include debug fields in the scan results"`
}

// Description returns an overview of this module.
func (flags *Redis) Description() string {
	return "Probe for Redis"
}

// Help returns the module's help string.
func (flags *Redis) Help() string {
	return ""
}

// Validate checks that the flags are valid.
func (flags *Redis) Execute(args []string) error {
	return nil
}
