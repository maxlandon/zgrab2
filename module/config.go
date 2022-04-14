package module

import (
	"net"
	"os"
)

// Config is the high level framework options that will be parsed
// from the command line.
type Config struct {
	OutputFileName     string `short:"o" long:"output-file" default:"-" description:"Output filename, use - for stdout"`
	InputFileName      string `short:"f" long:"input-file" default:"-" description:"Input filename, use - for stdin"`
	MetaFileName       string `short:"m" long:"metadata-file" default:"-" description:"Metadata filename, use - for stderr"`
	LogFileName        string `short:"l" long:"log-file" default:"-" description:"Log filename, use - for stderr"`
	LocalAddress       string `long:"source-ip" description:"Local source IP address to use for making connections"`
	Senders            int    `short:"s" long:"senders" default:"1000" description:"Number of send goroutines to use"`
	Debug              bool   `long:"debug" description:"Include debug fields in the output."`
	Flush              bool   `long:"flush" description:"Flush after each line of output."`
	GOMAXPROCS         int    `long:"gomaxprocs" default:"0" description:"Set GOMAXPROCS"`
	ConnectionsPerHost int    `long:"connections-per-host" default:"1" description:"Number of times to connect to each host (results in more output)"`
	ReadLimitPerHost   int    `long:"read-limit-per-host" default:"96" description:"Maximum total kilobytes to read for a single host (default 96kb)"`
	Prometheus         string `long:"prometheus" description:"Address to use for Prometheus server (e.g. localhost:8080). If empty, Prometheus is disabled."`
	// Multiple           MultipleCommand `command:"multiple" description:"Multiple module actions"`
	inputFile  *os.File
	outputFile *os.File
	metaFile   *os.File
	logFile    *os.File
	// inputTargets  InputTargetsFunc
	// outputResults OutputResultsFunc
	localAddr *net.TCPAddr
}
