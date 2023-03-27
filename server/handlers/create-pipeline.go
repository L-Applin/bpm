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
	io.WriteString(w, request.Host)
}
