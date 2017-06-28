package main

import (
	"os"
	"fmt"
	"strings"
)


const (
	error_v = iota
	warning_v
	info_v
	debug_v
)

var loggingLevelEnvVar = "GLLVM_OUTPUT_LEVEL"
var loggingFileEnvVar = "GLLVM_OUTPUT_FILE"


var loggingLevels = map[string]int{
	"ERROR":    error_v,
	"WARNING":  warning_v,
	"INFO":     info_v,
	"DEBUG":    debug_v,
}


var level = 0

var filePointer = os.Stderr


func init(){
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
	return func(format string, a ...interface{}){
		if level >= lvl {
			msg := fmt.Sprintf(format, a...)
			if !strings.HasSuffix(msg, "\n") {
				msg += "\n"
			}
			filePointer.WriteString(msg)
		}
	}
}

var logDebug   = makeLogger(debug_v)
var logInfo    = makeLogger(info_v)
var logWarning = makeLogger(warning_v)
var logError   = makeLogger(error_v)
