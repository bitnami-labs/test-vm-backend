// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"net/http"
	"strings"
)

func prepareCORSHeaders(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
}

func preflightReply(w http.ResponseWriter, r *http.Request, methods []string) {
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
	if requestedHeaders := r.Header.Get("Access-Control-Request-Headers"); requestedHeaders != "" {
		w.Header().Set("Access-Control-Allow-Headers", requestedHeaders)
	}
	w.WriteHeader(http.StatusNoContent)
}
