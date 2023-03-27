package main

import (
	"bpm/server/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/create-pipeline", handlers.CreatePipelineHandler{})

	err := http.ListenAndServe(":3333", nil)

	if err != nil {
		fmt.Errorf("error: %v", err.Error())
	}
}
