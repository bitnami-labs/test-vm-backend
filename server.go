package main

import (
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"
)

// VMServer is a http.Handler of VM REST requests
type VMServer struct {
	vmm     Cloud
	address string
}

// APIDoc dumps the API simple doc
func (s *VMServer) APIDoc() string {
	return `API
GET /vms             -> VMs JSON          # list all VMs
PUT /vms/launch/{id} -> Check status code # launch VM by id
PUT /vms/stop/{id}   -> Check status code # a VM by id
GET /vms/{id}        -> VM JSON           # inspect a VM by id
DELETE /vms/{id}     -> Check status code # delete a VM by id`
}

// ServeVM dispatchs the request to the correct method follwing the API schema
func (s *VMServer) ServeVM(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("<- %v %v\n", r.Method, r.URL.Path)
	switch {
	case matches(r, http.MethodPut, "/vms/launch/\\d+"):
		s.launch(w, r)
	case matches(r, http.MethodPut, "/vms/stop/\\d+"):
		s.stop(w, r)
	case matches(r, http.MethodGet, "/vms/\\d+"):
		s.inspect(w, r)
	case matches(r, http.MethodDelete, "/vms/\\d+"):
		s.delete(w, r)
	case matches(r, http.MethodGet, "/vms[/]?"):
		s.list(w, r)
	default:
		msg := fmt.Sprintf("%v %v not allowed", r.Method, r.URL.Path)
		http.Error(w, msg, http.StatusMethodNotAllowed)
	}
}

func matches(r *http.Request, method, pathRegex string) bool {
	if r.Method != method {
		return false
	}

	pattern := fmt.Sprintf("^%v$", pathRegex)
	matches, err := regexp.Match(pattern, []byte(r.URL.Path))
	dieOnError(err, "Error maching %q as %q", r.URL.Path, pattern)
	return matches
}

func (s *VMServer) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%v not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte(s.vmm.List().String()))
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
	w.Write([]byte(s.vmm.Inspect(id).String()))
}
