package bin

import (
	"github.com/zmap/zgrab2"
	// The module directory contains all flag configurations
	// for the default zgrab2 scan module, to be exposed to CLI
	"github.com/zmap/zgrab2/module"

	// All subdirectories are each protocol scanning code.
	"github.com/zmap/zgrab2/module/bacnet"
	"github.com/zmap/zgrab2/module/banner"
	"github.com/zmap/zgrab2/module/dnp3"
	"github.com/zmap/zgrab2/module/fox"
	"github.com/zmap/zgrab2/module/ftp"
	"github.com/zmap/zgrab2/module/http"
	"github.com/zmap/zgrab2/module/imap"
	"github.com/zmap/zgrab2/module/ipp"
	"github.com/zmap/zgrab2/module/jarm"
	"github.com/zmap/zgrab2/module/modbus"
	"github.com/zmap/zgrab2/module/mongodb"
	"github.com/zmap/zgrab2/module/mssql"
	"github.com/zmap/zgrab2/module/mysql"
	"github.com/zmap/zgrab2/module/ntp"
	"github.com/zmap/zgrab2/module/oracle"
	"github.com/zmap/zgrab2/module/pop3"
	"github.com/zmap/zgrab2/module/postgres"
	"github.com/zmap/zgrab2/module/redis"
	"github.com/zmap/zgrab2/module/siemens"
	"github.com/zmap/zgrab2/module/smb"
	"github.com/zmap/zgrab2/module/smtp"
	"github.com/zmap/zgrab2/module/telnet"
	"github.com/zmap/zgrab2/module/tls"
)

// The main code package has access to both module and scanners
// lists, so it takes care of populating each list, and then leaves
// the callers the choice of importing either or both lists.
var add = zgrab2.Register

// Each call below will:
// 1 - Make the scan module flags available as a CLI command.
// 2 - Register the scan code counterpart in Zgrab's main code.
func init() {
	// Name / Default port / Pointer to module flags / pointer to module scanner
	// -----------------------------------------------
	add("bacnet", 0xBAC0, &module.Bacnet{}, &bacnet.Scanner{})
	add("banner", 80, &module.Banner{}, &banner.Scanner{})
	add("dnp3", 20000, &module.DNP3{}, &dnp3.Scanner{})
	add("fox", 1911, &module.Fox{}, &fox.Scanner{})
	add("ftp", 21, &module.FTP{}, &ftp.Scanner{})
	add("http", 80, &module.HTTP{}, &http.Scanner{})
	add("imap", 143, &module.IMAP{}, &imap.Scanner{})
	add("ipp", 631, &module.IPP{}, &ipp.Scanner{})
	add("jarm", 443, &module.JARM{}, &jarm.Scanner{})
	add("modbus", 502, &module.Modbus{}, &modbus.Scanner{})
	add("mongodb", 27017, &module.MongoDB{}, &mongodb.Scanner{})
	add("mssql", 1433, &module.MSSQL{}, &mssql.Scanner{})
	add("mysql", 3306, &module.MySQL{}, &mysql.Scanner{})
	add("ntp", 123, &module.NTP{}, &ntp.Scanner{})
	add("oracle", 1521, &module.Oracle{}, &oracle.Scanner{})
	add("pop3", 110, &module.POP3{}, &pop3.Scanner{})
	add("postgres", 5432, &module.Postgres{}, &postgres.Scanner{})
	add("redis", 6379, &module.Redis{}, &redis.Scanner{})
	add("siemens", 102, &module.Siemens{}, &siemens.Scanner{})
	add("smb", 445, &module.SMB{}, &smb.Scanner{})
	add("smtp", 25, &module.SMTP{}, &smtp.Scanner{})
	add("ssh", 22, &module.SSH{}, &smtp.Scanner{}) // TODO: some default values to be set in ssh.go
	add("telnet", 23, &module.Telnet{}, &telnet.Scanner{})
	add("tls", 443, &module.TLS{}, &tls.Scanner{})
}
