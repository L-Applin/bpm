// Copyright (c) 2023 Olivier Lepage-Applin. All rights reserved.

package handlers

import (
	"bpm/log"
	"bpm/server/db"
	"bpm/utils"
	"encoding/json"
	"net/http"
	"strings"
)

type ProjectListResponse struct {
	Project []ProjectListElem `json:"projects"`
}

type ProjectListElem struct {
	Name      string   `json:"name"`
	Pipelines []string `json:"pipelines"`
}

func NewProjectsHandler() http.Handler {
	return New(ProjectsHandler{})
}

type ProjectsHandler struct {
}

func (h ProjectsHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	url := request.URL
	if strings.HasSuffix(url.String(), "/api/projects/") {
		listProjects(w, request)
	} else {
		singleProject(w, request)
	}
}

func singleProject(w http.ResponseWriter, request *http.Request) {
	splits := strings.Split(request.URL.String(), "/")
	projectName := splits[len(splits)-1]
	log.Debugf("command: Single project '%s'", projectName)

	infoJson := ProjectListElem{
		Name:      projectName,
		Pipelines: []string{"pipeline-1", "pipeline-2", "pipeline-3"},
	}
	b, err := json.Marshal(infoJson)
	if err != nil {
		log.Errorf("error while marshalling json: %$v", infoJson)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	log.Debugf("sending response: %s", string(b))
	if _, err := w.Write(b); err != nil {
		log.Errorf("error while writing to response: %#v", string(b))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug("single projects done")

}

func listProjects(w http.ResponseWriter, request *http.Request) {
	log.Debug("command: list projects")
	infos := db.GetProjectInfo()
	infoJson := utils.MapList[db.ProjectInfo, ProjectListElem](infos, func(info db.ProjectInfo) ProjectListElem {
		return ProjectListElem{
			Name:      info.Name,
			Pipelines: info.Pipelines,
		}
	})
	b, err := json.Marshal(ProjectListResponse{Project: infoJson})
	if err != nil {
		log.Errorf("error while marshaling to json: %#v", infoJson)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	log.Debugf("sending response: %s", string(b))
	if _, err := w.Write(b); err != nil {
		log.Errorf("error while writing to response: %#v", string(b))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug("list projects done")
}
