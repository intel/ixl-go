// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package log provides internal used log function like Debug and Info.
// User can change the log level by setting environment variable IXL_LOG_LEVEL.
package log

import (
	"log"
	"os"
)

var level = 0

const (
	lvlSilent = 0
	lvlInfo   = 1
	lvlDebug  = 2
)

func init() {
	if lvl := os.Getenv("IXL_LOG_LEVEL"); lvl != "" {
		switch lvl {
		case "SILENT":
			level = lvlSilent
		case "INFO":
			level = lvlInfo
		case "DEBUG":
			level = lvlDebug
		}
	}
}

// Debug prints debug level logs
func Debug(format string, args ...interface{}) {
	if level >= lvlDebug {
		log.Printf("[ixl-go]DEBG"+format, args...)
	}
}

// Info prints info level logs
func Info(format string, args ...interface{}) {
	if level >= lvlDebug {
		log.Printf("[ixl-go]INFO "+format, args...)
	}
}
