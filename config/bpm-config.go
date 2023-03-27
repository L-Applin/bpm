package config

import (
	"bpm/log"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	envPrefix       = "$"
	environmentsKey = "environments"
	envKey          = "env"
)

type PipelineConfiguration struct {
	Pipeline PipelineFile `yaml:"pipeline"`
}

type PipelineFile struct {
	File         string   `yaml:"file"`
	Description  string   `yaml:"description"`
	Project      string   `yaml:"project"`
	Environments []string `yaml:"environments"`
	PipelineFunc string   `yaml:"func"`
	Config       Config   `yaml:"configurations"`
}

type Config map[string]interface{}

func ParseConfigFile(file, env string) (PipelineConfiguration, error) {
	filename, _ := filepath.Abs(file)
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return PipelineConfiguration{}, fmt.Errorf("could not read config file '%s'\n", file)
	}
	c := PipelineConfiguration{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return PipelineConfiguration{}, fmt.Errorf("error unmarshalling config file '%s'\n", file)
	}
	found := false
	for _, e := range c.Pipeline.Environments {
		if e != env {
			delete(c.Pipeline.Config, envPrefix+e)
			continue
		}
		if envConf := c.Pipeline.Config[envPrefix+e]; envConf != nil {
			log.Debugf("found env '%s'\n", env)
			c.Pipeline.Config[envKey] = e
			delete(c.Pipeline.Config, environmentsKey)
			if confAsMap, ok := envConf.(Config); ok {
				c.Pipeline.Config = MergeConfigs(c.Pipeline.Config, confAsMap)
			}
			delete(c.Pipeline.Config, envPrefix+e)
		}
		found = true
	}
	if !found {
		return PipelineConfiguration{}, fmt.Errorf(
			"specified environment '%s' is not part of pipeline known environments: %s",
			env, c.Pipeline.Environments)
	}

	return c, nil
}

func MergeConfigs(initial, toMerge Config) Config {
	// if value for key is string, put it in the map
	// else if it is a map, if it does not exist
	config := Config{}
	for k, v := range initial {
		config[k] = v
	}
	for k, v := range toMerge {
		if _, ok := v.(string); ok {
			config[k] = v
		} else if c, ok := v.(Config); ok {
			if initial[k] == nil {
				initial[k] = Config{}
			}
			if initialDown, ok := initial[k].(Config); ok {
				config[k] = MergeConfigs(initialDown, c)
			} else {
				panic(fmt.Sprintf("%#v is not a config", initial[k]))
			}
		}
	}
	return config
}
