package zgrab2

import (
	"bytes"
	"errors"
	"os"

	"github.com/maxlandon/go-flags"
	"gopkg.in/ini.v1"

	zflags "github.com/zmap/zgrab2/modules/flags"
)

// MultipleCommand contains the command line options for running
type MultipleCommand zflags.MultipleCommand

// Execute - The multiple command is the only command that has a "non-empty" -thus fixed-
// command implementation.
func (x *MultipleCommand) Execute(args []string) error {
	if len(x.Args.Configs) == 0 {
		x.Args.Configs = append(x.Args.Configs, "-")
	}
	if x.ConfigFileName == config.InputFileName {
		return errors.New("cannot receive config file and input file from same source")
	}

	return nil
}

// Help returns a usage string that will be output at the command line
func (x *MultipleCommand) Help() string {
	return ""
}

// ParseScanners - The comand parses the supplied INI files, modifies any of
// its sections' names if there are multiple calls to the same module, and
// loads the contents of all files into a new flags.IniParser.
// @data is an optional INI configuration as bytes, such as those used in remote requests.
func (x *MultipleCommand) ParseScanners() (scanCommands []*flags.Command) {

	// The options that we must use when loading INI files
	iniOpts := ini.LoadOptions{
		AllowNonUniqueSections: true,
	}

	// The final INI command parser that will apply the file values
	// This parser already has default commands for each module.
	iniParser := flags.NewIniParser(parser)

	// The final INI configuration that we will use to parse
	// values from: in this one all sections have unique names,
	// and they might come from different files.
	var cfg = ini.Empty()

	// Temporary lists used to track which one needs new commands
	var unbound []string
	var registered []string

	// For each scan INI configuration file,
	// load and adjust sections if some are name-colliding
	for _, file := range x.Args.Configs {
		u, r := x.loadINI(file, cfg, iniOpts)
		unbound = append(unbound, u...)
		registered = append(registered, r...)
	}

	// We now have a single INI config with all our scans.
	// For each section, we register an appropriate command if needed
	x.addScanCommands(unbound, parser)

	// Both the INI configuration and the parser are set with
	// all specified scan modules, their flags and flags.Commands.
	// Run the parser to load all options.
	var cData bytes.Buffer
	cfg.WriteTo(&cData)
	iniParser.Parse(&cData)

	// And for each registered scan module command, add
	// to the list so that scanners can access their underlying data
	for _, scan := range registered {
		scanCommands = append(scanCommands, parser.Find(scan))
	}

	return
}

// loadINI - Given a path to an INI file, load it and adjust for any name-colliding sections.
func (x *MultipleCommand) loadINI(file string, cfg *ini.File, opts ini.LoadOptions) (u, r []string) {

	// Get the INI representation, or die tryin'
	var loaded *ini.File
	var err error
	// We either read from stdin
	if file == "-" {
		loaded, err = ini.LoadSources(opts, os.Stdin)

		// Or we are provided a remote request
	} else if file == "__remote_request__" && len(x.ConfigData) > 0 {
		loaded, err = ini.LoadSources(opts, x.ConfigData)

		// Or we must read a file somewhere
	} else {
		loaded, err = ini.LoadSources(opts, file)
	}
	if err != nil {
		return
	}

	// Cycle through sections
	for _, s := range loaded.Sections() {

		// The unique name for this scan section
		var name = s.Name()

		// If name found in registered modules,
		// add a nonce to the section title, and
		// mark this module as one needing a new command instance.
		if isRegisteredModule(s.Name(), r) {
			for isRegisteredModule(s.Name(), r) {
				name = name + "-"
			}
			u = append(u, name)
		}
		r = append(r, name)

		// Create a new INI section and copy
		// the contents of the old one into it.
		sec, _ := cfg.NewSection(name)
		sec.SetBody(s.Body())
		sec.Comment = s.Comment
	}

	return
}

// addScanCommands - Given a list of scan modules that have not yet their own command
// instance registered in the flags.Parser, trim name and add these commands.
func (x *MultipleCommand) addScanCommands(unbound []string, parser *flags.Parser) {
	for _, sec := range unbound {
		name := removeTrailingNonce(sec)
		mod := GetModule(name)

		// Descriptions for our new command
		var short, long string

		// Find an already registered command for the module
		cmd := parser.Find(name)
		if cmd == nil {
			short = name + " scanner"
			long = mod.Description()
		}
		short = cmd.ShortDescription
		long = cmd.LongDescription

		// And register our unique command to the parser
		parser.AddCommand(name, short, long, mod.NewFlags())
	}
}

func isRegisteredModule(mod string, reg []string) bool {
	for _, r := range reg {
		if mod == r {
			return true
		}
	}
	return false
}
