// Copyright (c) 2023 Olivier Lepage-Applin. All rights reserved.

package main

import (
	"bpm/log"
	"bpm/utils"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	CreatePipelineCmd       = "create-pipeline"
	ListProjectsInfoCommand = "projects"
)

// global arguments
var (
	bpmServer = flag.String("server", "http://localhost:3333", "BPM server url")
	verbose   = flag.Bool("verbose", false, "Enable verbose logging")
	v         = flag.Bool("v", false, "Enable verbose logging")
)

const BpmUsage = `BPM Usage
bpm [global flags] command [command flags] [command argument]
command:
  projects
  pipeline`

func main() {

	flag.Parse()
	if *verbose || *v {
		log.SetGlobalLogLevel(log.Levels.Debug)
	}
	log.Debugf("Global log level is: %s", log.GlobalLevel.Name)
	args := flag.Args()
	log.Debugf("args: %v", args)
	if len(args) < 1 {
		fmt.Println(BpmUsage)
		return
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "pipeline":
		pipeline(args)
	case "projects":
		projects(args)
	}
}

func pipeline(args []string) {
	subCmd, args := args[0], args[1:]
	switch subCmd {
	case "create":
		createPipeline(args)
	}

}

// bpm projects <project name>
func projects(args []string) {
	var uri string
	if len(args) == 0 {
		// GET <host>/api/projects
		uri = strings.Join([]string{bpmHost(), "api", "projects"}, "/")
	} else {
		projectName := args[0]
		// GET <host>/api/projects/<project name>
		uri = strings.Join([]string{bpmHost(), "api", "projects", projectName}, "/")
	}
	log.Debugf("sending request: GET %s", uri)
	res, err := http.Get(uri)
	if err != nil {
		ReportError(err)
		return
	}
	log.Debugf("response: %#v", res)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		ReportError(err)
		return
	}
	log.Debugf("body: %s", string(body))
}

type CreatePipelineArgs struct {
	file    string
	project string
}

func createPipeline(args []string) {
	cmdArgs := CreatePipelineArgs{}
	flagCmd := flag.NewFlagSet("create", flag.ExitOnError)
	flagCmd.StringVar(&cmdArgs.file, "file", "", "Pipeline yaml file to upload")
	flagCmd.StringVar(&cmdArgs.file, "f", "", "Pipeline yaml file to upload")
	flagCmd.StringVar(&cmdArgs.file, "project", "", "Project in which the pipeline will live")
	flagCmd.StringVar(&cmdArgs.file, "p", "", "Project in which the pipeline will live")
	if err := flagCmd.Parse(args); err != nil {
		ReportError(err)
		return
	}
	if err := DoCreatePipeline(cmdArgs); err != nil {
		ReportError(err)
	}
}

func DoCreatePipeline(args CreatePipelineArgs) error {
	uri := strings.Join([]string{bpmHost(), args.project, CreatePipelineCmd}, "/")
	f, err := os.Open(args.file)
	if err != nil {
		return err
	}
	base64FileReader := utils.Base64reader(f)
	resp, err := http.Post(uri, " application/octet-stream", base64FileReader)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err

	}
	fmt.Printf("%s", string(data))
	return nil
}

func bpmHost() string {
	return *bpmServer
}
