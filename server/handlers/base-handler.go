// Copyright (c) 2023 Olivier Lepage-Applin. All rights reserved.

package handlers

import (
	"bpm/log"
	"net/http"
)

func New(handler http.Handler) http.Handler {
	return BaseHandler{
		Handler: handler,
	}
}

// BaseHandler is the shared base entry point for handler all request to the BPM server
type BaseHandler struct {
	Handler http.Handler
}

// ServeHTTP delegates the handling of the request to the delegate handler
func (b BaseHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	log.Debugf("request received: %#v", request)
	b.Handler.ServeHTTP(w, request)
}
