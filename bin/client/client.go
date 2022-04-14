//go:build client
// +build client

package client

// import (
//         "fmt"
//         "log"
//         "os"
// )
//
// // Execute is the client implementation of our program. It must only parse
// // the values onto our various module/configuration structs, and then either
// // let the server-side work (either in the same binary, or with a remote invocation).
// func (r *root) Execute(args []string) (err error) {
//         fmt.Println("went client")
//         // Stats
//         // StartStats()
//         // defer DumpStats()
//
//         // First execute the cobra command line interface, for
//         // parsing our command line onto our module flags/commands
//         err = Zgrab.ParseFlags(os.Args[1:])
//         if err != nil {
//                 log.Fatal(err)
//         }
//
//         // And send our configuration remotely for execution.
//
//         // We need to get the name of the active command.
//         // The parser gives it to us, as the active one.
//         // command, err := ParseArgsAlt(args)
//         // command, err := zgrab2.ParseArgs(args)
//         // Blanked arg is positional arguments
//         // if err != nil {
//         //         // Outputting help is returned as an error. Exit successfully on help output.
//         //         flagsErr, ok := err.(*flags.Error)
//         //         if ok && flagsErr.Type == flags.ErrHelp {
//         //                 return
//         //         }
//         //         // Didn't output help. Unknown parsing error.
//         //         log.Fatalf("could not parse flags: %s", err)
//         // }
//
//         // Whether or not we have one or multiple
//         // scans to run (therefore multiple commands)
//         // we load them in a single list.
//         // By default, we have our CLI-invoked command.
//         // scanCommands := []*cobra.Command{command}
//         // var scanCommands = []*flags.Command{command}
//
//         // If multiple modules to run, set them up
//         // switch command.Data().(type) {
//         // case *zgrab2.MultipleCommand:
//         //         mult := command.Data().(*zgrab2.MultipleCommand)
//         //         scanCommands = append(scanCommands, mult.ParseScanners()...)
//         // }
//
//         // For each selected scan command, load the appropriate scanner
//         // for _, cmd := range scanCommands {
//         //         mod := GetModule(cmd.Name())
//         //         s := mod.NewScanner()
//         //         flags, _ := mod.NewFlags().(ScanFlags)
//         //         // flags, _ := cmd.Data().(zgrab2.ScanFlags)
//         //         s.Init(flags)
//         //         RegisterScan(cmd.Name(), s)
//         // }
//
//         // TODO: return the error to the caller instead of logging it here.
//         return
// }
