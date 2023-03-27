package examples

import (
	. "bpm"
	"bpm/agents"
	"bpm/agents/java"
	"bpm/config"
	"fmt"
)

func ReleasePipeline(config config.PipelineConfiguration) Pipeline {
	if config.Pipeline.Config["env"] == "prod" {
		fmt.Println("~ Generating PROD pipeline ~")
	}
	conf := config.Pipeline.Config
	buildAppPhase := NewPhase(
		"Build App - {{env}}",
		Phase{
			Output:   "/result-{{env}}",
			Agent:    java.Agent(java.Jdk11, java.Maven),
			Trigger:  GitWebHook("my-web-hook-url"),
			Callback: TestCallback,
			Steps: Steps{
				Script{
					Name: "Dummy Step - {{env}}",
					File: "examples/scripts/dummy-{{env}}.yml",
				},
				Step{
					Name: "Clone Repository - {{env}}",
					Commands: Commands{
						Command{"cd", "{{work-dir}}/{{git.dir}}"},
						Command{"git", "clone {{git.repo}}"},
					},
					Env: Env{
						"SOME_ENV_VAR": "{{demo}}",
					},
				},
				Step{
					Name: "Unit tests - {{env}}",
					Commands: Commands{
						Command{"./mvnw", "clean test -DskipTests {{maven-args}}"},
					},
				},
				Step{
					Name: "Maven Build - {{env}}",
					Commands: Commands{
						Command{"rm", "-r result"},
						Command{"./mvnw", "install -DskipTests {{maven-args}}"},
						Command{"mv", "target /result-{{env}}"},
					},
				},
				Step{
					Name: "Push to S3 - {{env}}",
					Commands: Commands{
						// dummy command
						Command{"aws", "s3 PutObject build.jar"},
					},
				},
			},
		})
	printResultPhase := NewPhase(
		"Print Result",
		Phase{
			Trigger: OnSuccess{
				AllOf: Dependencies{buildAppPhase.Id},
			},
			Input: "/result-{{env}}",
			Agent: agents.Bash,
			Steps: Steps{
				Parallel{
					Name: "Print result - {{env}}",
					Steps: Steps{
						Step{
							Name: "Print to std-out - {{env}}",
							Commands: Commands{
								Command{"ls", "-al /target"},
							},
						},
						Step{
							Name: "Print to socket - {{env}}",
							Commands: Commands{
								Command{"ls", "-al /target | nc localhost 3000"},
							},
							IgnoreFailure: true,
						},
					},
				},
				Step{
					Name: "Print Success - {{env}}",
					Commands: Commands{
						Command{"echo", "Success!!"},
					},
				},
			},
		})

	deployPhase := NewPhase("Deploy - {{env}}",
		Phase{
			Description: "test description for env: '{{env}}'",
			Trigger: OnSuccess{
				AllOf: Dependencies{buildAppPhase.Id, printResultPhase.Id},
			},
			Input:    "/result-{{env}}",
			Agent:    agents.Bash,
			Callback: DeployCallback,
			Steps: Steps{
				Step{
					Name: fmt.Sprintf("Deploy to AWS Lambda - %s", conf["env"]),
					Commands: Commands{
						Command{"aws", "lambda update-function-code --function-name dummy --zip-file fileb://dummy.zip"},
					},
				}},
		})
	return Pipeline{
		Name: "Release Pipeline",
		Phases: Phases{
			buildAppPhase,
			printResultPhase,
			deployPhase,
		}}
}

func DeployCallback(res *PhaseExecutionResult) {
	if res.Successful {
		fmt.Println("Successfull stage completed!")
	} else {
		fmt.Printf("Failed to complete phase '%s'\n", res.PhaseName)
	}
}

func TestCallback(res *PhaseExecutionResult) {
	if res.Successful {
		fmt.Println("Successfull stage completed!")
	} else {
		fmt.Printf("Failed to complete phase '%s'\n", res.PhaseName)
	}
}
