package handlers

import (
	"io"
	"net/http"
)

func NewCreatePipelineHandler() http.Handler {
	return CreatePipelineHandler{}
}

type CreatePipelineHandler struct {
}

func (CreatePipelineHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if b != nil {
		w.WriteHeader()
	}
}
