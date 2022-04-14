package ssh

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/lib/ssh"
	"github.com/zmap/zgrab2/module"
)

type SSHScanner struct {
	config *module.SSH
}

func (s *SSHScanner) Init(flags module.Scan) error {
	f, _ := flags.(*module.SSH)
	s.config = f
	return nil
}

func (s *SSHScanner) InitPerSender(senderID int) error {
	return nil
}

func (s *SSHScanner) GetName() string {
	return s.config.Name
}

func (s *SSHScanner) GetTrigger() string {
	return s.config.Trigger
}

func (s *SSHScanner) Scan(t zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	data := new(ssh.HandshakeLog)

	var port uint
	// If the port is supplied in ScanTarget, let that override the cmdline option
	if t.Port != nil {
		port = *t.Port
	} else {
		port = s.config.Port
	}
	portStr := strconv.FormatUint(uint64(port), 10)
	rhost := net.JoinHostPort(t.Host(), portStr)

	sshConfig := ssh.MakeSSHConfig()
	sshConfig.Timeout = s.config.Timeout
	sshConfig.ConnLog = data
	sshConfig.ClientVersion = s.config.ClientID
	sshConfig.HelloOnly = s.config.HelloOnly
	if err := sshConfig.SetHostKeyAlgorithms(s.config.HostKeyAlgorithms); err != nil {
		log.Fatal(err)
	}
	if err := sshConfig.SetKexAlgorithms(s.config.KexAlgorithms); err != nil {
		log.Fatal(err)
	}
	if err := sshConfig.SetCiphers(s.config.Ciphers); err != nil {
		log.Fatal(err)
	}
	sshConfig.Verbose = s.config.Verbose
	sshConfig.DontAuthenticate = s.config.CollectUserAuth
	sshConfig.GexMinBits = s.config.GexMinBits
	sshConfig.GexMaxBits = s.config.GexMaxBits
	sshConfig.GexPreferredBits = s.config.GexPreferredBits
	sshConfig.BannerCallback = func(banner string) error {
		data.Banner = strings.TrimSpace(banner)
		return nil
	}
	_, err := ssh.Dial("tcp", rhost, sshConfig)
	// TODO FIXME: Distinguish error types
	status := zgrab2.TryGetScanStatus(err)
	return status, data, err
}

// Protocol returns the protocol identifier for the scanner.
func (s *SSHScanner) Protocol() string {
	return "ssh"
}
