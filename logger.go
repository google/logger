/*
Copyright 2016 Google Inc. All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package logger offers simple cross platform logging for Windows and Linux.
// Available logging endpoints are event log (Windows), syslog (Linux), and 
// an io.Writer.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var (
	infoLog     *log.Logger
	errorLog    *log.Logger
	fatalLog    *log.Logger
	logLock     sync.Mutex
	initialized bool
)

// Init sets up logging and should be called before log functions, usually in
// the callers main(). Log functions can be called before Init(), but log
// output will only go to stderr (along with a warning).
func Init(name string, verbose, systemLog bool, logFile io.Writer) {
	var il, el io.Writer
	if systemLog {
		var err error
		il, el, err = setup(name)
		if err != nil {
			log.Fatal(err)
		}
	}

	iLogs := []io.Writer{logFile}
	eLogs := []io.Writer{logFile, os.Stderr}
	if verbose {
		iLogs = append(iLogs, os.Stdout)
	}
	if il != nil {
		iLogs = append(iLogs, il)
	}
	if el != nil {
		eLogs = append(eLogs, el)
	}

	flags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	infoLog = log.New(io.MultiWriter(iLogs...), "INFO: ", flags)
	errorLog = log.New(io.MultiWriter(eLogs...), "ERROR: ", flags)
	fatalLog = log.New(io.MultiWriter(eLogs...), "FATAL: ", flags)
	initialized = true
}

type severity int

const (
	sInfo = iota
	sError
	sFatal
)

func output(s severity, txt string) {
	logLock.Lock()
	defer logLock.Unlock()
	initText := "ERROR: Logging before logger.Init."
	switch s {
	case sInfo:
		if !initialized {
			fmt.Fprintf(os.Stderr, "%s\nINFO: %s\n", initText, txt)
			return
		}
		infoLog.Output(3, txt)
	case sError:
		if !initialized {
			fmt.Fprintf(os.Stderr, "%s\nERROR: %s\n", initText, txt)
			return
		}
		errorLog.Output(3, txt)
	case sFatal:
		if !initialized {
			fmt.Fprintf(os.Stderr, "%s\nFATAL: %s\n", initText, txt)
			return
		}
		fatalLog.Output(3, txt)
	default:
		panic(fmt.Sprintln("unrecognized severity:", s))
	}
}

// Info logs with the INFO severity.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	output(sInfo, fmt.Sprint(v...))
}

// Infoln logs with the INFO severity.
// Arguments are handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	output(sInfo, fmt.Sprintln(v...))
}

// Infof logs with the INFO severity.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	output(sInfo, fmt.Sprintf(format, v...))
}

// Error logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	output(sError, fmt.Sprint(v...))
}

// Errorln logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	output(sError, fmt.Sprintln(v...))
}

// Errorf logs with the Error severity.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	output(sError, fmt.Sprintf(format, v...))
}

// Fatal logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Print.
func Fatal(v ...interface{}) {
	output(sFatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalln logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Println.
func Fatalln(v ...interface{}) {
	output(sFatal, fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
func Fatalf(format string, v ...interface{}) {
	output(sFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}
