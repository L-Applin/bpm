// Copyright (c) 2023 Olivier Lepage-Applin. All rights reserved.

package log

import (
	"fmt"
	"strings"
	"time"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	grey  = "\033[90m"
	reset = "\033[0m"
)

type Color string

type Level struct {
	Rank  int
	Name  string
	Color Color
}

var Levels = struct {
	Error Level
	Info  Level
	Debug Level
	Trace Level
}{
	Error: Level{Rank: 0, Name: "ERROR", Color: red},
	Info:  Level{Rank: 1, Name: "INFO", Color: green},
	Debug: Level{Rank: 2, Name: "DEBUG", Color: grey},
	Trace: Level{Rank: 3, Name: "TRACE", Color: reset},
}

var (
	timeLayout = "2006-01-02T15:04:05"
)

var GlobalLevel = Levels.Info

func Log(level Level, log string) {
	if requiresLogging(level) {
		fmt.Printf("%s [%s%s%s] - %s\n",
			time.Now().Format(timeLayout),
			level.Color, level.Name, reset,
			log)
	}
}

func Logf(level Level, log string, a ...any) {
	if requiresLogging(level) {
		fmt.Printf("%s [%s%s%s] - %s\n",
			time.Now().Format(timeLayout),
			level.Color, level.Name, reset,
			fmt.Sprintf(log, a...))
	}
}

func Error(log string) {
	Log(Levels.Error, log)
}

func ErrorE(err error) {
	Log(Levels.Error, err.Error())
}

func Errorf(log string, a ...any) {
	Logf(Levels.Error, log, a)
}

func Info(log string) {
	Log(Levels.Info, log)
}

func Infof(log string, a ...any) {
	Logf(Levels.Info, log, a...)
}

func Debug(log string) {
	Log(Levels.Debug, log)
}

func Debugf(log string, a ...any) {
	Logf(Levels.Debug, log, a...)
}

func requiresLogging(level Level) bool {
	return !(level.Compare(GlobalLevel) > 0)
}

func (l Level) Compare(other Level) int {
	if l.Rank > other.Rank {
		return 1
	}
	return -1
}

func LevelFromString(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "info":
		return Levels.Info, nil
	case "debug":
		return Levels.Debug, nil
	case "trace":
		return Levels.Trace, nil
	default:
		return Levels.Info, fmt.Errorf("not a valid log level: '%s'", level)
	}
}

func SetGlobalLogLevel(level Level) {
	GlobalLevel = level
}

func SetGlobalLogLevelFromString(level string) {
	if l, err := LevelFromString(level); err == nil {
		SetGlobalLogLevel(l)
	}
}
