package log

import (
	"fmt"
	"time"
)

const (
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
	Info  Level
	Debug Level
	Trace Level
}{
	Info:  Level{Rank: 0, Name: "INFO", Color: green},
	Debug: Level{Rank: 1, Name: "DEBUG", Color: grey},
	Trace: Level{Rank: 1, Name: "TRACE", Color: reset},
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

func Info(log string) {
	Log(Levels.Info, log)
}

func Infof(log string, a ...any) {
	Logf(Levels.Info, log, a...)
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

func LevelFromString(level string) Level {
	switch level {
	case "info":
		return Levels.Info
	case "debug":
		return Levels.Debug
	case "trace":
		return Levels.Trace
	default:
		return Levels.Info
	}
}
