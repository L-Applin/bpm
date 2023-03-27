package bpm

import (
	"bpm/config"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"runtime"
)

type Command []string
type Commands []Command

type IStep interface {
	Preprocessor[IStep]
	Runner
	name() string
}

type Env map[string]string

// Step a Step is the basic unit of work. It is a list of commands to run
type Step struct {
	Name          string
	Commands      Commands
	IgnoreFailure bool
	Env           map[string]string
}

func (s Step) name() string {
	return s.Name
}

type Parallel struct {
	Name  string `json:"name"`
	Steps Steps  `json:"steps"`
}

func (ps Parallel) name() string {
	return ps.Name
}

// Phase a Phase is a collection of step that encapsulate a logical colection of work to be completed.
// Importantly, a Phase have their own Context: (inputs, outputs, Agent, etc).
type Phase struct {
	Name          string
	Description   Description
	Id            uuid.UUID // filled by bpm
	Agent         Agent
	Steps         Steps
	Input, Output string
	Trigger       Trigger
	Callback      Callback
}

type PhaseExecutionResult struct {
	PhaseName  string
	Successful bool
}

type Callback func(executionResult *PhaseExecutionResult)
type Description string
type Agent string
type Steps []IStep

type Script struct {
	Name string
	File string
}

func (fs Script) name() string {
	return fs.Name
}

type Trigger interface {
	Preprocessor[Trigger]
	GetTrigger() any
}

type Phases []Phase

// Dependencies are a list of other Phase ID which needs to succueed before triggering this task
type Dependencies []uuid.UUID

func NewPhase(name string, phase Phase) Phase {
	phase.Id = uuid.New()
	phase.Name = name
	return phase
}

type Pipeline struct {
	Config config.PipelineConfiguration
	Name   string
	Phases Phases
}

type Args struct {
	LogLevel        string
	ConfigFile      string
	Env             string
	AllowMissingVar bool
}

type WebHook struct {
	WebHook string
}

func GitWebHook(hook string) WebHook {
	return WebHook{WebHook: hook}
}

func (wh WebHook) GetTrigger() any {
	return wh.WebHook
}

type OnDeploy struct {
	Trigger bool
}

func (od OnDeploy) GetTrigger() any {
	return od
}

// OnPhasesSuccess requires all phases to pass successfully
func OnPhasesSuccess(phases ...Phase) onPhaseSuccess {
	var ids []uuid.UUID
	for _, x := range phases {
		ids = append(ids, x.Id)
	}
	return onPhaseSuccess{Requires: ids}
}

type onPhaseSuccess struct {
	Requires Dependencies
}

func (s onPhaseSuccess) GetTrigger() any {
	return s.Requires
}

type OnSuccess struct {
	AllOf Dependencies
	AnyOf Dependencies
}

func (os OnSuccess) GetTrigger() any {
	return os
}

type Cron struct {
	Pattern string
}

func (c Cron) GetTrigger() any {
	return c.Pattern
}

type OnError struct {
	phase uuid.UUID
}

func (oe OnError) GetTrigger() any {
	return oe.phase
}

func NewPhaseId() uuid.UUID {
	return uuid.New()
}

func Name(c Callback) string {
	return runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name()
}

func (c *Callback) MarshalJSON() ([]byte, error) {
	n := Name(*c)
	if n == "" {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", Name(*c))), nil
}

func (s Step) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name          string   `json:"name"`
		Commands      Commands `json:"commands"`
		IgnoreFailure bool     `json:"ignoreFailure"`
		Type          string   `json:"type"`
	}{
		Type:          "command",
		Name:          s.Name,
		Commands:      s.Commands,
		IgnoreFailure: s.IgnoreFailure,
	})
}

func (p Parallel) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  string `json:"name"`
		Steps Steps  `json:"steps"`
		Type  string `json:"type"`
	}{
		Type:  "parallel",
		Name:  p.Name,
		Steps: p.Steps,
	})
}

func (s Script) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name string `json:"name"`
		File string `json:"file"`
		Type string `json:"type"`
	}{
		Type: "script",
		Name: s.Name,
		File: s.File,
	})
}
