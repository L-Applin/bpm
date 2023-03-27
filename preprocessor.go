package bpm

import (
	"bpm/config"
	"fmt"
	"github.com/cbroglie/mustache"
)

type Preprocessor[T interface{}] interface {
	preprocess(config config.Config) T
}

func Preprocess(pipeline *Pipeline, args Args) {
	mustache.AllowMissingVariables = args.AllowMissingVar
	conf := pipeline.Config.Pipeline.Config
	for i := range pipeline.Phases {
		// Name
		processField(&pipeline.Phases[i].Name, conf)

		// Description
		pipeline.Phases[i].Description = Description(Process(string(pipeline.Phases[i].Description), conf))

		// Agent
		pipeline.Phases[i].Agent = Agent(Process(string(pipeline.Phases[i].Agent), conf))

		// Input
		processField(&pipeline.Phases[i].Input, conf)

		// Output
		processField(&pipeline.Phases[i].Output, conf)

		// Step
		processSteps(pipeline, i, conf)

		// Trigger
		processTrigger(pipeline, i, conf)
	}
}

func Process(data string, conf config.Config) string {
	rendered, err := mustache.Render(data, conf)
	if err != nil {
		panic(fmt.Errorf("error while processing '%s': %v", data, err))
	}
	return rendered
}

func processTrigger(pipeline *Pipeline, i int, conf config.Config) {
	trigger := pipeline.Phases[i].Trigger
	pipeline.Phases[i].Trigger = trigger.preprocess(conf)
}

func processSteps(pipeline *Pipeline, i int, conf config.Config) {
	for j, e := range pipeline.Phases[i].Steps {
		pipeline.Phases[i].Steps[j] = e.preprocess(conf)
	}
}

func processField(field *string, conf config.Config) {
	rendered, err := mustache.Render(*field, conf)
	if err != nil {
		panic(fmt.Errorf("error while processing field '%s': %v", *field, err))
	}
	*field = rendered
}

func (s Step) preprocess(conf config.Config) IStep {
	step := s
	processedName, err := mustache.Render(s.Name, conf)
	if err != nil {
		panic(err)
	}
	step.Name = processedName
	cmds := Commands{}
	for _, cmd := range s.Commands {
		processedC := Command{}
		for _, c := range cmd {
			proccessedCmd, err := mustache.Render(c, conf)
			if err != nil {
				panic(err)
			}
			processedC = append(processedC, proccessedCmd)
		}
		cmds = append(cmds, processedC)
	}
	step.Commands = cmds
	return step
}

func (ps Parallel) preprocess(conf config.Config) IStep {
	newStep := Parallel{}
	newName, err := mustache.Render(ps.Name, conf)
	if err != nil {
		panic(err)
	}
	newStep.Name = newName
	for _, s := range ps.Steps {
		newStep.Steps = append(newStep.Steps, s.preprocess(conf))
	}
	return newStep
}

func (wh WebHook) preprocess(conf config.Config) Trigger {
	rendered, err := mustache.Render(wh.WebHook, conf)
	if err != nil {
		panic(err)
	}
	return WebHook{
		WebHook: rendered,
	}
}

func (s onPhaseSuccess) preprocess(conf config.Config) Trigger {
	return s // do nothing
}

func (os OnSuccess) preprocess(conf config.Config) Trigger {
	return os // do nothing
}

func (c Cron) preprocess(conf config.Config) Trigger {
	rendered, err := mustache.Render(c.Pattern, conf)
	if err != nil {
		panic(err)
	}
	return Cron{
		Pattern: rendered,
	}
}

func (od OnDeploy) preprocess(config config.Config) Trigger {
	return od
}

func (fs Script) preprocess(conf config.Config) IStep {
	newStep := Script{}
	processedName, err := mustache.Render(fs.Name, conf)
	if err != nil {
		panic(err)
	}
	newStep.Name = processedName
	processedFile, err := mustache.Render(fs.File, conf)
	if err != nil {
		panic(err)
	}
	newStep.File = processedFile
	return newStep
}
