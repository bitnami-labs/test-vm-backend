package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// loadVMs loads the VM list from a JSON file (VMS_JSON)
func loadVMs() VMList {
	_, err := os.Stat(VMsJSON)
	if errors.Is(err, os.ErrNotExist) {
		saveVMs(defaultVMList)
	} else {
		dieOnError(err, "Error stating %q", VMsJSON)
	}
	f, err := os.Open(VMsJSON)
	dieOnError(err, "Error opening %q", VMsJSON)

	defer f.Close()
	vmsJSON, err := ioutil.ReadAll(f)
	dieOnError(err, "Error reading from %q", VMsJSON)

	vms := make(VMList, 0)
	err = json.Unmarshal(vmsJSON, &vms)
	dieOnError(err, "Error JSON-parsing from %q", VMsJSON)

	return vms
}

// saveVMs saves the VM list to a JSON file (VMS_JSON)
func saveVMs(vms VMList) {
	vmsJSON, err := json.Marshal(vms)
	dieOnError(err, "Error writing JSON for %q", VMsJSON)

	err = ioutil.WriteFile(VMsJSON, vmsJSON, 0644)
	dieOnError(err, "Error saving %q", VMsJSON)
}

func main() {
	vms := VMServer{Cloud{vms: loadVMs()}, ":8080"}

	fmt.Printf("Server listening at %v\n", vms.address)
	fmt.Println(vms.APIDoc())
	http.HandleFunc("/", vms.ServeVM)

	dieOnError(http.ListenAndServe(vms.address, nil), "Error listening on %v", vms.address)
}
