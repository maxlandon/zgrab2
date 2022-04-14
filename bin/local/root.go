package local

import (

	// Commands and struct parsers.
	"github.com/octago/sflags/gen/gcobra"
	"github.com/spf13/cobra"

	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/bin/client"
	"github.com/zmap/zgrab2/module"

	// Include the scanners as an anonymous import:
	// The init function will add both the client (flags)
	// and the server (scan) parts of each module.
	_ "github.com/zmap/zgrab2/bin"
)

// Zgrab is the root command for the zgrab2 binary.
// It is created at init time and does:
// - Scans the root configuration flags and creates the command.
// - Scans all available module flags and creates their commands.
var Zgrab = gcobra.Parse(zgrab2.GetConfig())

func init() {
	// 1 - Fetch a list of scan modules generated as cobra commands,
	// with everything set, but without any PreRun/Run/RunE implementation.
	// This gives us the client side of zgrab2.
	commands := client.ModuleCommands()

	// 2 - We bind the same run implementation to all these scan commands,
	// since they will end up sending their scans to the same place.
	for _, command := range commands {
		command.RunE = runScans
	}

	// 3 - Bind the commands to the root zgrab2 command.
	// The scanners have been registered through an init
	// function in the package imported anonymously above.
	Zgrab.AddCommand(commands...)
}

// runScans is the only run command used by zgrab2, since all scans
// go through the same registration and execution process. The way
// cobra works will also directly find the appropriate command.
func runScans(cmd *cobra.Command, args []string) (err error) {
	// If the multiple command was actually called, the next few
	// steps/function calls will have no effect, since the multiple
	// command has registered scans on its own. The party goes on
	// again when zgrab2.RunScans is called below.
	//
	// Otherwise, we find here a single match for the command,
	// which we pass to the server package for init and run.
	// TODO: move this, avoids calling the module package. Since
	// zgrab2 has access to it, it can find the associated flags by name.
	mod := module.GetModule(cmd.Name())

	// Since we are in a local binary, we directly call the server
	// package for registering our scan data which, again, might be
	// nil if the invoked command was "multiple".
	zgrab2.RegisterScan(cmd.Name(), mod)

	// Finally, start all the scans, blocking until they all terminate
	// or an error is raised, in case returned to the caller.
	return zgrab2.RunScans()
}

// Execute is the "root" local implementation of zgrab2. In other words,
// this takes care both of parsing command-line flags and registering/
// running the scans associated. This file/function has thus access to
// the entire zgrab2 library code, as opposed to the client-only root.
// func (r *root) Execute(args []string) (err error) {
//         // Get the actual subcommand that was invoked.
//         // All commands match with a given scan module, except
//         // the "multiple" one which is actually a subcommand.
//         cmd, _, err := Zgrab.Find(os.Args[1:])
//         if err != nil {
//                 log.Fatal(err)
//         }
//
//         // If the multiple command was actually called, the next few
//         // steps/function calls will have no effect, since the multiple
//         // command has registered scans on its own. The party goes on
//         // again when zgrab2.RunScans is called below.
//         //
//         // Otherwise, we find here a single match for the command,
//         // which we pass to the server package for init and run.
//         module := modules.GetModule(cmd.Name())
//
//         // Since we are in a local binary, we directly call the server
//         // package for registering our scan data which, again, might be
//         // nil if the invoked command was "multiple".
//         zgrab2.RegisterScan(cmd.Name(), module)
//
//         // Finally, start all the scans, blocking until they all terminate.
//         err = zgrab2.RunScans()
//         if err != nil {
//                 log.Fatal(err)
//         }
//
//         return
// }
