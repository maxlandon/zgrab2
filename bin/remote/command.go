package remote

import "github.com/maxlandon/gonsole"

// private interface used to represent a parser that can add commands to itself.
type parser interface {
	AddCommand(string, string, string, string, []string, gonsole.Commander) *gonsole.Command
}

// RemoteScan is a top level command that can be registered
// to an external command parser, for remote scan command exec.
type RemoteScan struct{}

// Execute is the scan implementation for remote execution.
func (s *RemoteScan) Execute(args []string) (err error) {
	return
}

// AddRemoteScanCommands binds all remote equivalents of Zgrab
// module commands, so that you can run scans/sets remotely with
// the same parsing and CLI functionality precision. There are
// also methods available in the library so that you can process
// the output of the scans, which will thus come back from remote.
func AddRemoteScanCommands(parser parser) (err error) {

	// Root zgrab command
	parser.AddCommand("zgrab",
		"All Zgrab scan module commands", "",
		"scans",
		[]string{""},
		&RemoteScan{},
	)

	// All remote equivalents of the module commands

	// Multiple Scans
	parser.AddCommand("multiple",
		"Run multiple scans, from one or more INI files on disk, or from stdin (-)", "",
		"multi",
		[]string{""},
		&MultipleCommand{},
	)
	// multi.AddArgumentCompletion("Configs", gonsole.Command)

	return
}
