package zgrab2

// var parser *flags.Parser
//
// func init() {
//         parser = flags.NewParser(&config, flags.Default)
// }
//
// // NewIniParser creates and returns a ini parser initialized
// // with the default parser.
// func NewIniParser() *flags.IniParser {
//         return flags.NewIniParser(parser)
// }
//
// // AddGroup exposes the parser's AddGroup function, allowing extension
// // of the global arguments.
// func AddGroup(shortDescription string, longDescription string, data interface{}) {
//         parser.AddGroup(shortDescription, longDescription, data)
// }
//
// // AddCommand adds a module to the parser and returns a pointer to
// // a flags.command object or an error.
// func AddCommand(command string, shortDescription string, longDescription string, port int, m ScanModule) (*flags.Command, error) {
//         cmd, err := parser.AddCommand(command, shortDescription, longDescription, m.NewFlags())
//         if err != nil {
//                 return nil, err
//         }
//         cmd.FindOptionByLongName("port").Default = []string{strconv.FormatUint(uint64(port), 10)}
//         cmd.FindOptionByLongName("name").Default = []string{command}
//         modules[command] = m
//         return cmd, nil
// }
//
// // ParseArgs - Returns the active command of the parser, that is,
// // the command that the user entered, like `zgrab2 ssh`.
// func ParseArgs(args []string) (cmd *flags.Command, err error) {
//         _, err = parser.ParseArgs(args)
//         if err == nil {
//                 ValidateFrameworkConfiguration()
//         }
//         return parser.Active, err
// }
