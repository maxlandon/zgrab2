package flags

// MultipleCommand contains the command line options for running.
type MultipleCommand struct {
	Args struct {
		Configs []string `description:"path to one or more config files, use - for stdin"`
	} `positional-args:"yes"`
	ConfigFileName  string `short:"c" long:"config-file" default:"-" description:"Config filename, use - for stdin"`
	ContinueOnError bool   `long:"continue-on-error" description:"If proceeding protocols error, do not run following protocols (default: true)"`
	BreakOnSuccess  bool   `long:"break-on-success" description:"If proceeding protocols succeed, do not run following protocols (default: false)"`

	// remote
	ConfigData []byte
}

// NewRemoteMultiple - Directly pass an INI file as bytes and return a multiple command.
func NewRemoteMultiple(data []byte) *MultipleCommand {
	x := &MultipleCommand{
		ConfigData: data,
	}
	x.Args.Configs = append(x.Args.Configs, "__remote_request__")
	return x
}
