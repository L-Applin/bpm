package main

import (
	"bpm"
	"bpm/config"
	"bpm/log"
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"plugin"
)

func main() {
	args := bpm.Args{}
	CliArgs(&args)
	conf, err := config.ParseConfigFile(args.ConfigFile, args.Env)
	if err != nil {
		panic(err)
	}
	pipeline, err := getPipeline(conf)
	if err != nil {
		panic(err)
	}
	pipeline.Config = conf
	bpm.Preprocess(&pipeline, args)
	log.Debugf("%s\n", asJson(&pipeline))
	if err := runPipeline(&pipeline); err != nil {
		panic(err)
	}

}

func getPipeline(conf config.PipelineConfiguration) (bpm.Pipeline, error) {
	fileName := conf.Pipeline.File
	outName := fileName[:len(fileName)-len(filepath.Ext(fileName))] + ".so"
	log.Debugf("Compiling pipeline to plugin '%s'", outName)
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outName, fileName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return bpm.Pipeline{}, fmt.Errorf("error while loading pipeline: %v", string(output))
	}
	log.Debugf("Compiling pipeline to plugin '%s'", outName)
	plug, _ := plugin.Open(outName)
	pip, _ := plug.Lookup(conf.Pipeline.PipelineFunc)
	loaded, _ := pip.(func(config.PipelineConfiguration) bpm.Pipeline)
	return loaded(conf), nil
}

func runPipeline(p *bpm.Pipeline) error {
	log.Infof("Running pipeline '%s'", p.Name)
	for _, phase := range p.Phases {
		if err := phase.Run(bpm.Context{
			Phase: phase,
		}); err != nil {
			return fmt.Errorf("error while running pipeline: %s", err)
		}
	}
	return nil
}

func asJson(o any) string {
	empJSON, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(empJSON)
}

func CliArgs(args *bpm.Args) {
	flag.StringVar(&args.ConfigFile, "config", "", "path to the config file to use")
	flag.StringVar(&args.Env, "env", "default", "the environment to use")
	flag.BoolVar(&args.AllowMissingVar, "allowMissingVar", false,
		"If missing preprocessed variable should fail pre-processing")
	flag.StringVar(&args.LogLevel, "log", "Info", "Log level")
	flag.Parse()
	log.SetGlobalLogLevelFromString(args.LogLevel)
	log.Debugf("%#v\n", *args)
}
