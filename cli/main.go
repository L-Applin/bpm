package main

import (
	"flag"
	"os"
)

const (
	CreatePipelineCmd = "create-pipeline"
)

func main() {
	CliFlags()

	switch os.Args[1] {
	case CreatePipelineCmd:

	}
}

func CliFlags() {
	CreatePipelineCommand()
}

func CreatePipelineCommand() {
	flagCmd := flag.NewFlagSet("create-pipeline", flag.ExitOnError)
	flagCmd.String("file", "", "Pipeline yaml file to upload")
}
