// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// VMServer is a http.Handler of VM REST requests
type VMServer struct {
	vmm Cloud
}

type serverHandler func(s *VMServer, w http.ResponseWriter, r *http.Request)

type idHandlerFunc func(id int, w http.ResponseWriter, r *http.Request)

func mustCompileAnchored(pattern string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
}

// MethodSpec defines a single method on an entrypoint
type MethodSpec struct {
	Method   string
	BodySpec string
	Doc      string
	Handler  serverHandler
}

// EndpointSpec defined a path endpoint and all its methods
type EndpointSpec struct {
	DisplayPath string
	Path        *regexp.Regexp
	Methods     []MethodSpec
}

// APISpec specifies endpoint paths and their implemented methods
var APISpec = []EndpointSpec{
	{
		DisplayPath: "/vms",
		Path:        mustCompileAnchored(`/vms[/]?`),
		Methods: []MethodSpec{
			{
				http.MethodGet, "VMs JSON", "list All VMs",
				func(s *VMServer, w http.ResponseWriter, r *http.Request) {
					s.list(w, r)
				},
			},
		},
	},
	{
		DisplayPath: "/vms/{vm_id}/launch",
		Path:        mustCompileAnchored(`/vms/\d+/launch[/]?`),
		Methods: []MethodSpec{
			{
				http.MethodPut, "", "launch VM by id",
				func(s *VMServer, w http.ResponseWriter, r *http.Request) {
					s.requestIDfor(s.launch, 2, w, r)
				},
			},
		},
	},
	{
		DisplayPath: "/vms/{vm_id}/stop",
		Path:        mustCompileAnchored(`/vms/\d+/stop[/]?`),
		Methods: []MethodSpec{
			{
				http.MethodPut, "", "stop VM by id",
				func(s *VMServer, w http.ResponseWriter, r *http.Request) {
					s.requestIDfor(s.stop, 2, w, r)
				},
			},
		},
	},
	{
		DisplayPath: "/vms/{vm_id}",
		Path:        mustCompileAnchored(`/vms/\d+`),
		Methods: []MethodSpec{
			{
				http.MethodGet, "VM JSON", "inspect a VM by id",
				func(s *VMServer, w http.ResponseWriter, r *http.Request) {
					s.requestIDfor(s.inspect, 2, w, r)
				},
			},
			{
				http.MethodDelete, "", "delete a VM by id",
				func(s *VMServer, w http.ResponseWriter, r *http.Request) {
					s.requestIDfor(s.delete, 2, w, r)
				},
			},
		},
	},
}

// NewVMServer returns a new VM server
func NewVMServer(vms VMs) *VMServer {
	return &VMServer{Cloud{vms: vms}}
}

// WriteAPIDoc dumps the API simple doc onto the given writer
func (s *VMServer) WriteAPIDoc(w io.Writer) {
	fmt.Fprintln(w, "API:")
	for _, endpoint := range APISpec {
		for _, m := range endpoint.Methods {
			bodySpec := m.BodySpec
			if bodySpec == "" {
				bodySpec = "Check status code"
			}
			fmt.Fprintf(w, "%v\t%-20v\t-> %-20v\t# %v\n",
				m.Method, endpoint.DisplayPath, bodySpec, m.Doc)
		}
	}
}

// ServeHTTP dispatchs the request to the correct method follwing the API schema
func (s *VMServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("<- %v %v", r.Method, r.URL.Path)
	for _, endpoint := range APISpec {
		if endpoint.Path.MatchString(r.URL.Path) {
			for _, m := range endpoint.Methods {
				if r.Method == m.Method {
					m.Handler(s, w, r)
					return
				}
			}
		}
	}
	msg := fmt.Sprintf("%v %v not allowed", r.Method, r.URL.Path)
	http.Error(w, msg, http.StatusMethodNotAllowed)
}

func matches(r *http.Request, method string, pathRegex *regexp.Regexp) bool {
	if r.Method != method {
		return false
	}
	return pathRegex.MatchString(r.URL.Path)
}

func (s *VMServer) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%v not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, s.vmm.List().String())
}

func (s *VMServer) requestIDfor(f idHandlerFunc, pos int, w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path.Base(pathParts[pos]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f(id, w, r)
}

func (s *VMServer) launch(id int, w http.ResponseWriter, r *http.Request) {
	if _, err := s.vmm.Launch(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (s *VMServer) stop(id int, w http.ResponseWriter, r *http.Request) {
	if _, err := s.vmm.Stop(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (s *VMServer) delete(id int, w http.ResponseWriter, r *http.Request) {
	if err := s.vmm.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
}

func (s *VMServer) inspect(id int, w http.ResponseWriter, r *http.Request) {
	vm, _ := s.vmm.Inspect(id)
	if _, err := fmt.Fprint(w, vm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
