package bin

import (
	"os"

	"github.com/maxlandon/go-flags"
	log "github.com/sirupsen/logrus"

	"github.com/zmap/zgrab2"
)

// ZGrab2Main should be called by func main() in a binary. The caller is
// responsible for importing any modules in use. This allows clients to easily
// include custom sets of scan modules by creating new main packages with custom
// sets of ZGrab modules imported with side-effects.
func ZGrab2Main() {
	ZGrab2Args(os.Args[1:])
}

// ZGrab2Args can be called with an arbitrary list of command args. The caller is
// responsible for importing any modules in use. This allows clients to easily
// include custom sets of scan modules by creating new main packages with custom
// sets of ZGrab modules imported with side-effects.
func ZGrab2Args(args []string) {

	// Stats
	zgrab2.StartStats()
	defer zgrab2.DumpStats()

	// We need to get the name of the active command.
	// The parser gives it to us, as the active one.
	command, err := zgrab2.ParseArgs(args)

	// Blanked arg is positional arguments
	if err != nil {
		// Outputting help is returned as an error. Exit successfuly on help output.
		flagsErr, ok := err.(*flags.Error)
		if ok && flagsErr.Type == flags.ErrHelp {
			return
		}
		// Didn't output help. Unknown parsing error.
		log.Fatalf("could not parse flags: %s", err)
	}

	// Whether or not we have one or multiple
	// scans to run (therefore multiple commands)
	// we load them in a single list.
	// By default, we have our CLI-invoked command.
	var scanCommands = []*flags.Command{command}

	// If multiple modules to run, set them up
	switch command.Data().(type) {
	case *zgrab2.MultipleCommand:
		mult := command.Data().(*zgrab2.MultipleCommand)
		scanCommands = append(scanCommands, mult.ParseScanners()...)
	}

	// For each selected scan command, load the appropriate scanner
	for _, cmd := range scanCommands {
		mod := zgrab2.GetModule(cmd.Name)
		s := mod.NewScanner()
		flags, _ := cmd.Data().(zgrab2.ScanFlags)
		s.Init(flags)
		zgrab2.RegisterScan(cmd.Name, s)
	}

	// Finally, start all the scans, blocking until they all terminate.
	err = zgrab2.RunScans()
	if err != nil {
		log.Fatal(err)
	}
}
