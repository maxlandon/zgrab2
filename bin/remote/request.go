package remote

import "encoding/json"

// Request - A type passed between a client and a scan server
type Request struct {
	Name      string   // The name of the command to run
	CmdFlags  []byte   // The flags of the invoked command, as bytes
	ScanFiles [][]byte // A list of INI files as bytes
}

// Remote is a command with a special execution method, which
// forwards its content to a remote zgrab scan process.
type Remote struct {
	isRemote bool // Embedders notify whether or not to run this command.

	// Operating
	// Flags - The embedder will assign itself to this,
	// so that we can unmarshal it around its remote peer
	Flags interface{}
}

// Marshal - A client embedding the module commands in its own application should call this method
// its own implementation of the Execute(args []string) command. This allows you to prepare a fully
// ready buffer containing all the data this is needed to run the very same command, remotely.
func (r *Remote) Marshal(name string, flags interface{}, args []string) (data []byte, err error) {
	if !r.isRemote {
		return []byte{}, nil
	}

	// Marshal our data to JSON
	bData, err := json.Marshal(flags)
	if err != nil {
		return
	}

	// Wrap everything in a request and marshal
	req := Request{
		Name:     name,
		CmdFlags: bData,
	}
	return json.Marshal(req)
}
