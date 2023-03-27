package main

import (
	. "bpm"
	"bpm/agents"
	. "bpm/config"
)

func HelloWorld(config PipelineConfiguration) Pipeline {
	phase1 := NewPhase(
		"Print Hello World",
		Phase{
			Description: "Prints 'Hello World' using bash",
			Agent:       agents.Bash,
			Trigger:     OnDeploy{Trigger: true},
			Output:      "var/",
			Steps: Steps{
				Step{
					Name: "Hello World Step",
					Commands: Commands{
						Command{"echo", "{{to-say}} > var/out.txt"},
					},
				},
				Step{
					Name: "Done",
					Commands: Commands{
						Command{"echo", "Done!"},
						Command{"echo", "Exiting!"},
					},
				}},
		})
	phase2 := NewPhase(
		"In response",
		Phase{
			Description: "Triggered after phase1",
			Agent:       agents.Bash,
			Trigger:     OnSuccess{AllOf: Dependencies{phase1.Id}},
			Steps: Steps{
				Step{
					Name: "cat",
					Commands: Commands{
						Command{"ls", "-al", "var/"},
						Command{"cat", "var/out.txt"},
					},
				},
			},
		})
	return Pipeline{
		Name: "Hello World pipeline",
		Phases: Phases{
			phase1,
			phase2,
		}}
}
