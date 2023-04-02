package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
)

const (
	CreatePipelineCmd = "create-pipeline"
)

func main() {

	switch os.Args[1] {
	case CreatePipelineCmd:
		args := CreatePipelineCommand()
		CreatePipeline(args)
	}
}

type CreatePipelineArgs struct {
	file string
	project string
}

func CreatePipelineCommand() CreatePipelineArgs {
	args := CreatePipelineArgs{}
	flagCmd := flag.NewFlagSet("create-pipeline", flag.ExitOnError)
	flagCmd.StringVar(&args.file, "file", "", "Pipeline yaml file to upload")
	flagCmd.StringVar(&args.file, "f", "", "Pipeline yaml file to upload")
	flagCmd.StringVar(&args.file, "project", "", "Project in which the pipeline will live")
	flagCmd.StringVar(&args.file, "p", "", "Project in which the pipeline will live")
	return args
}

func CreatePipeline(args CreatePipelineArgs) {
	uri := strings.Join([]string{bpmHost(), args.project, CreatePipelineCmd}, "/")
	http.Post(
}

func bpmHost() string {
	return "localhost:3333/api"
}
