package main

import (
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// VMServer is a http.Handler of VM REST requests
type VMServer struct {
	vmm     Cloud
	address string
}

type serverHandler func(s *VMServer, w http.ResponseWriter, r *http.Request)

type idHandlerFunc func(id int, w http.ResponseWriter, r *http.Request)

func mustCompileAnchored(pattern string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%v$", pattern))
}

// APISpec is the data source for API docs and mappings
var APISpec = []struct {
	method   string
	path     *regexp.Regexp
	bodySpec string
	doc      string
	handler  serverHandler
}{
	{
		http.MethodGet, wrapNCompile("/vms[/]?"), "VMs JSON", "list All VMs",
		func(s *VMServer, w http.ResponseWriter, r *http.Request) { s.list(w, r) },
	},
	{
		http.MethodPut, wrapNCompile("/vms/launch/\\d+"), "", "launch VM by id",
		func(s *VMServer, w http.ResponseWriter, r *http.Request) { s.requestIDfor(s.launch, w, r) },
	},
	{
		http.MethodPut, wrapNCompile("/vms/stop/\\d+"), "", "stop a VM by id",
		func(s *VMServer, w http.ResponseWriter, r *http.Request) { s.requestIDfor(s.stop, w, r) },
	},
	{
		http.MethodGet, wrapNCompile("/vms/\\d+"), "VM JSON", "inspect a VM by id",
		func(s *VMServer, w http.ResponseWriter, r *http.Request) { s.requestIDfor(s.inspect, w, r) },
	},
	{
		http.MethodDelete, wrapNCompile("/vms/\\d+"), "VM JSON", "delete a VM by id",
		func(s *VMServer, w http.ResponseWriter, r *http.Request) { s.requestIDfor(s.delete, w, r) },
	},
}

// APIDoc dumps the API simple doc
func (s *VMServer) APIDoc() string {
	var sb strings.Builder
	sb.WriteString("API:\n")
	for _, endpoint := range APISpec {
		bodySpec := endpoint.bodySpec
		if bodySpec == "" {
			bodySpec = "Check status code"
		}
		sb.WriteString(fmt.Sprintf("%v\t%-20v\t-> %-20v\t# %v\n",
			endpoint.method, endpoint.path, bodySpec, endpoint.doc))
	}
	return sb.String()
}

// ServeVM dispatchs the request to the correct method follwing the API schema
func (s *VMServer) ServeVM(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("<- %v %v\n", r.Method, r.URL.Path)
	for _, endpoint := range APISpec {
		if matches(r, endpoint.method, endpoint.path) {
			endpoint.handler(s, w, r)
			return
		}
	}
	msg := fmt.Sprintf("%v %v not allowed", r.Method, r.URL.Path)
	http.Error(w, msg, http.StatusMethodNotAllowed)
}

func matches(r *http.Request, method string, pathRegex *regexp.Regexp) bool {
	if r.Method != method {
		return false
	}
	return pathRegex.Match([]byte(r.URL.Path))
}

func (s *VMServer) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%v not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte(s.vmm.List().String()))
}

func (s *VMServer) requestIDfor(f idHandlerFunc, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f(id, w, r)
}

func (s *VMServer) launch(id int, w http.ResponseWriter, r *http.Request) {
	if err := s.vmm.Launch(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (s *VMServer) stop(id int, w http.ResponseWriter, r *http.Request) {
	if err := s.vmm.Stop(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (s *VMServer) delete(id int, w http.ResponseWriter, r *http.Request) {
	s.vmm.Delete(id)
}

func (s *VMServer) inspect(id int, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s.vmm.Inspect(id).String()))
}
