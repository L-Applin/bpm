// Copyright (c) 2023 Olivier Lepage-Applin. All rights reserved.

package db

type BpmDb interface {
	ProjectsInfo() ([]ProjectInfo, error)
	ProjectInfo(projectName string) (ProjectInfo, error)
	SavePipeline(base64File []byte) error
}

type ProjectInfo struct {
	Name      string
	Pipelines []string
}
