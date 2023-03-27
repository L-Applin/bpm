# BPM
Automation system. You write pipeline in pure go code. You can hook phase listeners that you write in go
```yaml
pipeline:
  file: "hello-world.go"
  func: "examples.HelloWorld"
  description: |
    hello world example for bpm
  project: "Hello World Project"
  environments: [ "demo" ]
  configurations:
    to-say: |
      Text from the config file
```
```go
package examples

import (
    . "bpm"
    "bpm/agents"
    . "bpm/config"
)

func HelloWorld(config PipelineConfiguration) Pipeline {
    return Pipeline{
        Name: "Hello World pipeline",
        Phases: Phases{
            NewPhase("Print Hello World",
                Phase{
                    Description: "Prints 'Hello World' using bash",
                    Agent:       agents.Bash,
                    Trigger:     OnDeploy{Trigger: true},
                    Steps: Steps{
                        Step{
                            Name: "Hello World Step",
                            Commands: Commands{
                                Command{"echo", "{{to-say}}"},
                            },
                        },
                        Step{
                            Name: "Done",
                            Commands: Commands{
                                Command{"echo", "Done!"},
                                Command{"echo", "Exiting!"},
                            },
                        }},
                }),
        }}
}
```


# Pipeline syntax
## Pipeline
```go
package bpm
type Pipeline struct {
	Phases Phases
	Config Config
}
```
### Config
```go
package bpm
type Config map[string]string
```
Configs are the value that will be replaced in strings with 'mustache', ie `Input:   "/result-{{env}}",`
`{{env}}` wil be replaced by the value in the map associated with the key `env`.

### Phases
```go
package bpm
type Phase struct {
	Name          string
	Description   Description
	Id            uuid.UUID
	Agent         Agent
	Steps         Steps
	Input, Output string
	Trigger       Trigger
}
```
#### Name
The name of the phase, for human to read

#### Description
```go
package bpm
type Description string
```
A description for the phase, for documentation purpose

#### Id
A unique uuid used to cross-reference Phases

#### Agent
```go
package bpm
type Agent string
```
An agent is a type of environment in which the phase is being executed in.

#### Steps

#### Input

#### Output

#### Trigger