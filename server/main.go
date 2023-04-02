package main

import (
	"bpm/log"
	"bpm/server/handlers"
	"net/http"
)

func main() {
	http.Handle("/create-pipeline", handlers.CreatePipelineHandler{})

	err := http.ListenAndServe(":3333", nil)

	if err != nil {
		log.ErrorE(err)
	}
}
