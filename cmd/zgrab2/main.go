package main

import (
	"log"

	// Automatic completion generation.
	"github.com/octago/sflags/gen/gcomp"

	// Import paths of zgrab2 components:.
	// bin/local    => everything in one binary, with CLI.
	// bin/client   => only the modules CLI interface.
	// bin/server   => only the scanners, to be served remotely.
	"github.com/zmap/zgrab2/bin/local"

	"github.com/zmap/zgrab2/module"
)

// main is the entrypoint of a local scanner binary:
// It must parse the command line onto our modules in the
// client package, then actually run their associated scanners
// in the server package.
func main() {
	// Generate the completions for the entire zgrab2 application.
	gcomp.Generate(local.Zgrab, &module.Config{}, nil)

	// The root command will roughly follow these steps:
	// 1 - simply parses the command-line flags (no exec)
	// 2 - finds the command invoked (eg. ssh/http/multiple)
	// 3 - passes the struct flags to the corresponding scanner.
	// 4 - Runs the scanners that we have set up in 3).
	if err := local.Zgrab.Execute(); err != nil {
		log.Fatal(err)
	}
}
