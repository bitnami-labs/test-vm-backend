// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Version of the program.
// "Development" by default, but released binaries will have versions such as:
// '2020.09.10.0'
var Version = "Development"

// loadVMs loads the VM list from a JSON file (VMS_JSON)
func loadVMs() (VMs, error) {
	log.Printf("Loading fake Cloud state from local file %q", VMsJSON)
	_, err := os.Stat(VMsJSON)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Missing %q, generating one...", VMsJSON)
		if err := saveVMs(defaultVMs); err != nil {
			return nil, fmt.Errorf("error generating default %q: %v", VMsJSON, err)
		}
		log.Printf("Tip: You can tweak %q adding VMs or changing states for next run.", VMsJSON)
	} else if err != nil {
		return nil, fmt.Errorf("error stating %q: %v", VMsJSON, err)
	}
	f, err := os.Open(VMsJSON)
	if err != nil {
		return nil, fmt.Errorf("error opening %q: %v", VMsJSON, err)
	}

	defer f.Close()
	vmsJSON, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading %q: %v", VMsJSON, err)
	}

	vms := make(VMs, 0)
	err = json.Unmarshal(vmsJSON, &vms)
	if err != nil {
		return nil, fmt.Errorf("error JSON-parsing %q: %v", VMsJSON, err)
	}

	return vms, nil
}

// saveVMs saves the VM list to a JSON file (VMS_JSON)
func saveVMs(vms VMs) error {
	vmsJSON, err := json.Marshal(vms)
	if err != nil {
		return fmt.Errorf("error writing JSON for %q: %v", VMsJSON, err)
	}

	err = ioutil.WriteFile(VMsJSON, vmsJSON, 0644)
	if err != nil {
		return fmt.Errorf("error saving %q: %v", VMsJSON, err)
	}
	return nil
}

func printDefaultsTo(w io.Writer, fs *flag.FlagSet) {
	defer func(saved io.Writer) {
		fs.SetOutput(saved)
	}(fs.Output())
	fs.SetOutput(w)
	fs.PrintDefaults()
}

func mainE() error {
	log.Printf("Test-VMBackend version %s", Version)
	var address string
	flag.StringVar(&address, "address", ":8080", "Listen address for the backend")
	flag.Parse()
	vms, err := loadVMs()
	if err != nil {
		return fmt.Errorf("error loading VMs initial state: %v", err)
	}
	server := VMServer{Cloud{vms: vms}, address}

	log.Printf("Server listening at %v", server.address)
	server.WriteAPIDoc(os.Stdout)
	http.HandleFunc("/", server.ServeVM)
	err = http.ListenAndServe(server.address, nil)
	if err != nil && strings.Contains(err.Error(), "address already in use") {
		var sb strings.Builder
		fmt.Fprintln(&sb, err.Error())
		fmt.Fprintf(&sb, "^ You can avoid binding issues by using the address flag:\n")
		printDefaultsTo(&sb, flag.CommandLine)
		return fmt.Errorf(sb.String())
	}
	return err
}

func main() {
	if err := mainE(); err != nil {
		log.Fatal(err)
	}
}
