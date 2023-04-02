package handlers

import (
	"bpm/log"
	"encoding/base64"
	"io"
	"net/http"
)

func NewBpmHandler() http.Handler {
	return New(BpmHandler{})
}

type BpmHandler struct {
}

func (BpmHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	decoded, _ := base64.StdEncoding.DecodeString(string(b))
	log.Debugf("body: %s", decoded)
	w.Write(decoded)
}
