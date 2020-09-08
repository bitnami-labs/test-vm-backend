package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// loadVMs loads the VM list from a JSON file (VMS_JSON)
func loadVMs() (VMList, error) {
	fmt.Printf("Loading fake Cloud state from local file %q\n", VMsJSON)
	_, err := os.Stat(VMsJSON)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Missing %q, generating one...\n", VMsJSON)
		if err := saveVMs(defaultVMList); err != nil {
			return nil, fmt.Errorf("error generating default %q: %v", VMsJSON, err)
		}
		fmt.Printf("Tip: You can tweak %q adding VMs or changing states for next run.\n", VMsJSON)
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

	vms := make(VMList, 0)
	err = json.Unmarshal(vmsJSON, &vms)
	if err != nil {
		return nil, fmt.Errorf("error JSON-parsing %q: %v", VMsJSON, err)
	}

	return vms, nil
}

// saveVMs saves the VM list to a JSON file (VMS_JSON)
func saveVMs(vms VMList) error {
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

func mainE() error {
	vms, err := loadVMs()
	if err != nil {
		return fmt.Errorf("error loading VMs initial state: %v", err)
	}
	server := VMServer{Cloud{vms: vms}, ":8080"}

	fmt.Printf("Server listening at %v\n", server.address)
	fmt.Println(server.APIDoc())
	http.HandleFunc("/", server.ServeVM)
	return http.ListenAndServe(server.address, nil)
}

func main() {
	if err := mainE(); err != nil {
		log.Fatal(err)
	}
}
