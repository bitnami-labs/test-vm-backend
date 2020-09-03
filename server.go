package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
)

// VMServer is a http.Handler of VM REST requests
type VMServer struct {
	vmm     Cloud
	address string
}

func (s *VMServer) serveVM(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.inspect(w, r)
	case http.MethodDelete:
		s.delete(w, r)
	default:
		http.Error(w, fmt.Sprintf("%v not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func (s *VMServer) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%v not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	vms, err := json.Marshal(s.vmm.List())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(vms)
}

func (s *VMServer) launch(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.vmm.Launch(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (s *VMServer) stop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.vmm.Stop(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (s *VMServer) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.vmm.Delete(id)
}

func (s *VMServer) inspect(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vm, err := json.Marshal(s.vmm.Inspect(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(vm)
}

// loadVMs loads the VM list from a JSON file (VMS_JSON)
func loadVMs() VMList {
	_, err := os.Stat(VMsJSON)
	if err == os.ErrNotExist {
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
