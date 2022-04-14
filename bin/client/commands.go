package client

import (
	"strconv"

	"github.com/octago/sflags/gen/gcobra"
	"github.com/spf13/cobra"
	"github.com/zmap/zgrab2/module"
)

// ModuleCommands fetches a list of all scan modules available and translates
// them into a list of cobra commands, to be attached to a cobra root command.
//
// Note that all these commands have an empty (but non-nil) Run implementation.
// It is the responsibility of the caller to assign them a new Run, depending
// on the setup of the zgrab2 binary (client/local/server).
func ModuleCommands() (commands []*cobra.Command) {
	// In the end, this call returns nothing more than a list of
	// generic structs with some tags and an Execute() implementation.
	mods, ports := module.GetAll()

	for command, mod := range mods {
		port := ports[command]

		cmd := gcobra.Parse(mod)
		// TODO: remove this as well
		if cmd == nil {
			continue
			// return errors.New("failed to parse struct into cobra command")
		}

		// Adjust name and descriptions
		cmd.Use = command
		cmd.Short = mod.Description()
		cmd.Long = mod.Help()

		// Set default flags
		cmd.Flags().Lookup("port").DefValue = strconv.FormatUint(uint64(port), 10)
		cmd.Flags().Lookup("name").DefValue = command

		// All modules share the same command implementation,
		// since we are only interested in parsing their flags.
		// This has also the advantage of transparently binding
		// either a local implementation, or a client-only one,
		// in which scanning is actually performed remotely.
		// cmd.RunE = Zgrab.RunE

		// And add to the cobra root command
		commands = append(commands, cmd)
	}

	return commands
}
