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
	defaultLogger *Logger
	logLock       sync.Mutex
)

const (
	flags    = log.Ldate | log.Lmicroseconds | log.Lshortfile
	initText = "ERROR: Logging before logger.Init.\n"
)

func initialize() {
	defaultLogger = &Logger{
		infoLog:  log.New(os.Stderr, initText+"INFO: ", flags),
		errorLog: log.New(os.Stderr, initText+"ERROR: ", flags),
		fatalLog: log.New(os.Stderr, initText+"FATAL: ", flags),
	}
}

func init() {
	initialize()
}

// Init sets up logging and should be called before log functions, usually in
// the caller's main(). Default log functions can be called before Init(), but log
// output will only go to stderr (along with a warning).
// The first call to Init populates the default logger and returns the
// generated logger, subsequent calls to Init will only return the generated
// logger.
func Init(name string, verbose, systemLog bool, logFile io.Writer) *Logger {
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

	var l Logger
	l.infoLog = log.New(io.MultiWriter(iLogs...), "INFO: ", flags)
	l.errorLog = log.New(io.MultiWriter(eLogs...), "ERROR: ", flags)
	l.fatalLog = log.New(io.MultiWriter(eLogs...), "FATAL: ", flags)
	l.initialized = true

	logLock.Lock()
	defer logLock.Unlock()
	if !defaultLogger.initialized {
		defaultLogger = &l
	}

	return &l
}

type severity int

const (
	sInfo = iota
	sError
	sFatal
)

// A Logger represents an active logging object. Multiple loggers can be used
// simultaneously even if they are using the same same writers.
type Logger struct {
	infoLog     *log.Logger
	errorLog    *log.Logger
	fatalLog    *log.Logger
	initialized bool
}

func (l *Logger) output(s severity, txt string) {
	logLock.Lock()
	defer logLock.Unlock()
	switch s {
	case sInfo:
		l.infoLog.Output(3, txt)
	case sError:
		l.errorLog.Output(3, txt)
	case sFatal:
		l.fatalLog.Output(3, txt)
	default:
		panic(fmt.Sprintln("unrecognized severity:", s))
	}
}

// Info logs with the INFO severity.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	l.output(sInfo, fmt.Sprint(v...))
}

// Infoln logs with the INFO severity.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Infoln(v ...interface{}) {
	l.output(sInfo, fmt.Sprintln(v...))
}

// Infof logs with the INFO severity.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(sInfo, fmt.Sprintf(format, v...))
}

// Error logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	l.output(sError, fmt.Sprint(v...))
}

// Errorln logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Errorln(v ...interface{}) {
	l.output(sError, fmt.Sprintln(v...))
}

// Errorf logs with the Error severity.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(sError, fmt.Sprintf(format, v...))
}

// Fatal logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Fatal(v ...interface{}) {
	l.output(sFatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalln logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Fatalln(v ...interface{}) {
	l.output(sFatal, fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.output(sFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Info calls the default logger's Info.
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Infoln calls the default logger's Infoln.
func Infoln(v ...interface{}) {
	defaultLogger.Infoln(v...)
}

// Infof calls the default logger's Infof.
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Error calls the default logger's Error.
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Errorln calls the default logger's Errorln.
func Errorln(v ...interface{}) {
	defaultLogger.Errorln(v...)
}

// Errorf calls the default logger's Errorf.
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// Fatal calls the default logger's Fatal.
func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

// Fatalln calls the default logger's Fatalln.
func Fatalln(v ...interface{}) {
	defaultLogger.Fatalln(v...)
}

// Fatalf calls the default logger's Fatalln.
func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}
