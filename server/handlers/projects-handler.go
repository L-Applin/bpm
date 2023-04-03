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
	dbbpm, err := db.MockDefault()
	if err != nil {
		panic("cannot init mock db")
	}
	return New(ProjectsHandler{
		db: dbbpm,
	})
}

type ProjectsHandler struct {
	db db.BpmDb
}

func (h ProjectsHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if strings.HasSuffix(path, "/api/projects/") {
		listProjects(w, h.db)
	} else {
		splits := strings.Split(path, "/")
		prjName := splits[len(splits)-1]
		singleProject(w, prjName, h.db)
	}
}

func listProjects(w http.ResponseWriter, dbBpm db.BpmDb) {
	log.Debug("command: list projects")
	infos, err := dbBpm.ProjectsInfo()
	if err != nil {
		log.ErrorE(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
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

func singleProject(w http.ResponseWriter, projectName string, dbBpm db.BpmDb) {
	log.Debugf("command: project info '%s'", projectName)
	info, err := dbBpm.ProjectInfo(projectName)
	if err != nil {
		log.ErrorE(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	infoJson := ProjectListElem{
		Name:      info.Name,
		Pipelines: info.Pipelines,
	}
	b, err := json.Marshal(infoJson)
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
