// Copyright (c) 2023 Olivier Lepage-Applin. All rights reserved.

// store pipelines in the file system. Used only for development and testing.

package db

import (
	"bpm/log"
	"bpm/utils"
	"fmt"
	"os"
	"strings"
)

func MockDefault() (BpmDb, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	defaultLocation := strings.Join([]string{userHome, ".bpm"}, string(os.PathSeparator))
	return Mock(defaultLocation)
}

func Mock(root string) (BpmDb, error) {
	if os.MkdirAll(root, os.ModePerm) != nil {
		return nil, fmt.Errorf("unable to create directory '%s'\n", root)
	}
	return FileSystemDb{
		Root: root,
	}, nil
}

type FileSystemDb struct {
	Root string
}

func (fsdb FileSystemDb) ProjectsInfo() ([]ProjectInfo, error) {
	fileInfos, err := os.ReadDir(strings.Join([]string{fsdb.Root, "projects"}, string(os.PathSeparator)))
	if err != nil {
		return []ProjectInfo{}, err
	}
	var projectInfos []ProjectInfo
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			pipelineInfo, err := os.ReadDir(strings.Join([]string{fsdb.Root, "projects", fileInfo.Name()}, string(os.PathSeparator)))
			if err != nil {
				return []ProjectInfo{}, err
			}
			projectInfos = append(projectInfos, ProjectInfo{
				Name:      fileInfo.Name(),
				Pipelines: utils.MapList(pipelineInfo, func(d os.DirEntry) string { return d.Name() }),
			})
		}
	}
	return projectInfos, nil
}

func (fsdb FileSystemDb) ProjectInfo(projectName string) (ProjectInfo, error) {
	fileInfos, err := os.ReadDir(strings.Join([]string{fsdb.Root, "projects", projectName}, string(os.PathSeparator)))
	log.Debugf("file info: %#v", fileInfos)
	if err != nil {
		return ProjectInfo{}, err
	}
	pipelines := utils.MapList(fileInfos, func(d os.DirEntry) string { return d.Name() })
	return ProjectInfo{
		Name:      projectName,
		Pipelines: pipelines,
	}, nil
}

func (fsdb FileSystemDb) SavePipeline(base64File []byte) error {
	//TODO implement me
	panic("implement me")
}
