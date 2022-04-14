package module

type SSH struct {
	Base              `group:"base"`
	ClientID          string `long:"client" description:"Specify the client ID string to use" default:"SSH-2.0-Go"`
	KexAlgorithms     string `long:"kex-algorithms" description:"Set SSH Key Exchange Algorithms"`
	HostKeyAlgorithms string `long:"host-key-algorithms" description:"Set SSH Host Key Algorithms"`
	Ciphers           string `long:"ciphers" description:"A comma-separated list of which ciphers to offer."`
	CollectUserAuth   bool   `long:"userauth" description:"Use the 'none' authentication request to see what userauth methods are allowed"`
	GexMinBits        uint   `long:"gex-min-bits" description:"The minimum number of bits for the DH GEX prime." default:"1024"`
	GexMaxBits        uint   `long:"gex-max-bits" description:"The maximum number of bits for the DH GEX prime." default:"8192"`
	GexPreferredBits  uint   `long:"gex-preferred-bits" description:"The preferred number of bits for the DH GEX prime." default:"2048"`
	HelloOnly         bool   `long:"hello-only" description:"Limit scan to the initial hello message"`
	Verbose           bool   `long:"verbose" description:"Output additional information, including SSH client properties from the SSH handshake."`
}

// Description returns an overview of this module.
func (m *SSH) Description() string {
	return "Fetch an SSH server banner and collect key exchange information"
}

func (f *SSH) Help() string {
	return ""
}

func (f *SSH) Execute(args []string) error {
	return nil
}

// type SSHModule struct {
// }

// func init() {
//         var sshModule SSHModule
//         cmd, err := zgrab2.AddCommand("ssh", "SSH Banner Grab", sshModule.Description(), 22, &sshModule)
//         if err != nil {
//                 log.Fatal(err)
//         }
//         s := ssh.MakeSSHConfig() //dummy variable to get default for host key, kex algorithm, ciphers
//         cmd.FindOptionByLongName("host-key-algorithms").Default = []string{strings.Join(s.HostKeyAlgorithms, ",")}
//         cmd.FindOptionByLongName("kex-algorithms").Default = []string{strings.Join(s.KeyExchanges, ",")}
//         cmd.FindOptionByLongName("ciphers").Default = []string{strings.Join(s.Ciphers, ",")}
// }

// func (m *SSHModule) NewFlags() interface{} {
//         return new(SSHFlags)
// }
//
// func (m *SSHModule) NewScanner() zgrab2.Scanner {
//         return new(SSHScanner)
// }
//
