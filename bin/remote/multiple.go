package remote

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/ini.v1"

	"github.com/zmap/zgrab2/modules/flags"
)

// MultipleCommand is a client-side implementation that will
// package a bunch of INI files and send them a remote instance.
type MultipleCommand flags.MultipleCommand

// Execute is the implementation for remote multiple scanning
func (m *MultipleCommand) Execute(args []string) (err error) {

	// The final INI configuration that we will use to parse
	// values from: in this one all sections have unique names,
	// and they might come from different files.
	var cfg = ini.Empty()

	// The options that we must use when loading INI files
	opts := ini.LoadOptions{
		AllowNonUniqueSections: true,
	}

	// Return the marshalled version of the command,
	// containing all the INI data needed for scans.
	for _, file := range m.Args.Configs {
		if err = m.loadINI(file, cfg, opts); err != nil {
			return err
		}
	}

	// Written to the configuration bytes
	buffer := bytes.Buffer{}
	cfg.WriteTo(&buffer)
	m.ConfigData = buffer.Bytes()

	return
}

// Response is the implementation of a reflag.ClientCommand type.
func (m *MultipleCommand) Response(data []byte) (err error) {
	return
}

// loadINI - Given a path to an INI file, load it and adjust for any name-colliding sections.
func (m *MultipleCommand) loadINI(file string, cfg *ini.File, opts ini.LoadOptions) (err error) {

	// Get the INI representation, or die tryin'
	var loaded *ini.File
	if file == "-" {
		loaded, err = ini.LoadSources(opts, os.Stdin)
	} else {
		loaded, err = ini.LoadSources(opts, file)
	}
	if err != nil {
		return fmt.Errorf("Error loading INI: %s", err)
	}

	// Cycle through sections
	for _, s := range loaded.Sections() {

		// The unique name for this scan section
		var name = s.Name()

		// Create a new INI section and copy
		// the contents of the old one into it.
		sec, _ := cfg.NewSection(name)
		sec.SetBody(s.Body())
		sec.Comment = s.Comment
	}

	return
}
