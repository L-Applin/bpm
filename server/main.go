package main

import (
	"bpm/log"
	"bpm/server/handlers"
	"net/http"
)

func main() {
	log.SetGlobalLogLevel(log.Levels.Debug)
	log.Info("Starting to listen on port 3333")
	http.Handle("/api/", handlers.NewBpmHandler())
	http.Handle("/api/projects/", handlers.NewProjectsHandler())

	err := http.ListenAndServe(":3333", nil)

	if err != nil {
		log.ErrorE(err)
	}
}
