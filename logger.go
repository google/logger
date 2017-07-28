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
	defaultLogger *logger
	initialized   bool
)

func init() {
	defaultLogger = &logger{}
}

// Init sets up logging and should be called before log functions, usually in
// the callers main(). Default log functions can be called before Init(), but log
// output will only go to stderr (along with a warning).
// The first call to Init populates the default logger and returns the
// generated logger, subsequent calls to Init will only return the generated
// logger.
func Init(name string, verbose, systemLog bool, logFile io.Writer) *logger {
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

	var l logger
	l.infoLog = log.New(io.MultiWriter(iLogs...), "INFO: ", flags)
	l.errorLog = log.New(io.MultiWriter(eLogs...), "ERROR: ", flags)
	l.fatalLog = log.New(io.MultiWriter(eLogs...), "FATAL: ", flags)

	if !initialized {
		defaultLogger = &l
		initialized = true
	}

	return &l
}

type severity int

const (
	sInfo = iota
	sError
	sFatal
)

type logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger
	logLock  sync.Mutex
}

func (l *logger) output(s severity, txt string) {
	l.logLock.Lock()
	defer l.logLock.Unlock()
	initText := "ERROR: Logging before logger.Init."
	switch s {
	case sInfo:
		if !initialized {
			fmt.Fprintf(os.Stderr, "%s\nINFO: %s\n", initText, txt)
			return
		}
		l.infoLog.Output(3, txt)
	case sError:
		if !initialized {
			fmt.Fprintf(os.Stderr, "%s\nERROR: %s\n", initText, txt)
			return
		}
		l.errorLog.Output(3, txt)
	case sFatal:
		if !initialized {
			fmt.Fprintf(os.Stderr, "%s\nFATAL: %s\n", initText, txt)
			return
		}
		l.fatalLog.Output(3, txt)
	default:
		panic(fmt.Sprintln("unrecognized severity:", s))
	}
}

// Info logs with the INFO severity.
// Arguments are handled in the manner of fmt.Print.
func (l *logger) Info(v ...interface{}) {
	l.output(sInfo, fmt.Sprint(v...))
}

// Infoln logs with the INFO severity.
// Arguments are handled in the manner of fmt.Println.
func (l *logger) Infoln(v ...interface{}) {
	l.output(sInfo, fmt.Sprintln(v...))
}

// Infof logs with the INFO severity.
// Arguments are handled in the manner of fmt.Printf.
func (l *logger) Infof(format string, v ...interface{}) {
	l.output(sInfo, fmt.Sprintf(format, v...))
}

// Error logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Print.
func (l *logger) Error(v ...interface{}) {
	l.output(sError, fmt.Sprint(v...))
}

// Errorln logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Println.
func (l *logger) Errorln(v ...interface{}) {
	l.output(sError, fmt.Sprintln(v...))
}

// Errorf logs with the Error severity.
// Arguments are handled in the manner of fmt.Printf.
func (l *logger) Errorf(format string, v ...interface{}) {
	l.output(sError, fmt.Sprintf(format, v...))
}

// Fatal logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Print.
func (l *logger) Fatal(v ...interface{}) {
	l.output(sFatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalln logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Println.
func (l *logger) Fatalln(v ...interface{}) {
	l.output(sFatal, fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
func (l *logger) Fatalf(format string, v ...interface{}) {
	l.output(sFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Info calls the default loggers Info.
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Infoln calls the default loggers Infoln.
func Infoln(v ...interface{}) {
	defaultLogger.Infoln(v...)
}

// Infof calls the default loggers Infof.
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Error calls the default loggers Error.
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Errorln calls the default loggers Errorln.
func Errorln(v ...interface{}) {
	defaultLogger.Errorln(v...)
}

// Errorf calls the default loggers Errorf.
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// Fatal calls the default loggers Fatal.
func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

// Fatalln calls the default loggers Fatalln.
func Fatalln(v ...interface{}) {
	defaultLogger.Fatalln(v...)
}

// Fatalf calls the default loggers Fatalln.
func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}
