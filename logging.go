package main

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

const (
	loggingLevelEnvVar = "GLLVM_OUTPUT_LEVEL"
	loggingFileEnvVar  = "GLLVM_OUTPUT_FILE"
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
	if envLevelStr := os.Getenv(loggingLevelEnvVar); envLevelStr != "" {
		if envLevelVal, ok := loggingLevels[envLevelStr]; ok {
			level = envLevelVal
		}
	}
	if envFileStr := os.Getenv(loggingFileEnvVar); envFileStr != "" {
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

var logDebug = makeLogger(debugV)
var logInfo = makeLogger(infoV)
var logWarning = makeLogger(warningV)
var logError = makeLogger(errorV)

func logFatal(format string, a ...interface{}) {
	logError(format, a...)
	os.Exit(1)
}
