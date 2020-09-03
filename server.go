package main

import (
	"fmt"
	"net/http"
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
