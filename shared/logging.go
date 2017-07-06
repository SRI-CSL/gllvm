package shared

import (
	"fmt"
	"os"
	"strings"
)

const (
	errorV = iota
	warningV
	infoV
	debugV
)

var loggingLevels = map[string]int{
	"ERROR":   errorV,
	"WARNING": warningV,
	"INFO":    infoV,
	"DEBUG":   debugV,
}

var level = 0

var filePointer = os.Stderr

func init() {
	if envLevelStr := os.Getenv("GLLVM_OUTPUT_LEVEL"); envLevelStr != "" {
		if envLevelVal, ok := loggingLevels[envLevelStr]; ok {
			level = envLevelVal
		}
	}
	if envFileStr := os.Getenv("GLLVM_OUTPUT_FILE"); envFileStr != "" {
		if loggingFP, err := os.OpenFile(envFileStr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err == nil {
			filePointer = loggingFP
		}
	}
}

func makeLogger(lvl int) func(format string, a ...interface{}) {
	return func(format string, a ...interface{}) {
		if level >= lvl {
			msg := fmt.Sprintf(format, a...)
			if !strings.HasSuffix(msg, "\n") {
				msg += "\n"
			}
			filePointer.WriteString(msg)
		}
	}
}

//LogDebug logs to the configured stream if the logging level is DEBUG.
var LogDebug = makeLogger(debugV)
//LogInfo logs to the configured stream if the logging level is INFO or lower.
var LogInfo = makeLogger(infoV)
//LogWarning logs to the configured stream if the logging level is WARNING or lower.
var LogWarning = makeLogger(warningV)
//LogError logs to the configured stream if the logging level is ERROR or lower.
var LogError = makeLogger(errorV)

//LogFatal logs to the configured stream and then exits.
func LogFatal(format string, a ...interface{}) {
	LogError(format, a...)
	os.Exit(1)
}
