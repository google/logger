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
	"errors"
	"fmt"
	"io"
	"os"
)

type severityType string

const (
	sInfo  = severityType("Info")
	sError = severityType("Error")
	sFatal = severityType("Fatal")
)

type message struct {
	sev severityType
	msg string
}

func (m message) String() string {
	return fmt.Sprintf("[%v]: %v", m.sev, m.msg)
}

type logger struct {
	msg chan message
}

// NewLogger sets up logging - should be called before log receiver methods.
func NewLogger(name string, w io.Writer) (*logger, error) {
	if len(name) == 0 {
		return nil, errors.New("null logger name")
	}
	l := &logger{
		msg: make(chan message),
	}
	go l.run(w)
	return l, nil
}

func (l *logger) run(w io.Writer) {
	for {
		m := <-l.msg
		fmt.Fprintln(w, m.String())
	}
}

// Info logs with the INFO severity.
// Arguments are handled in the manner of fmt.Print.
func (l *logger) Info(v ...interface{}) {
	l.msg <- message{sInfo, fmt.Sprint(v...)}
}

// Infoln logs with the INFO severity.
// Arguments are handled in the manner of fmt.Println.
func (l *logger) Infoln(v ...interface{}) {
	l.msg <- message{sInfo, fmt.Sprintln(v...)}
}

// Infof logs with the INFO severity.
// Arguments are handled in the manner of fmt.Printf.
func (l *logger) Infof(format string, v ...interface{}) {
	l.msg <- message{sInfo, fmt.Sprintf(format, v...)}
}

// Error logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Print.
func (l *logger) Error(v ...interface{}) {
	l.msg <- message{sError, fmt.Sprint(v...)}
}

// Errorln logs with the ERROR severity.
// Arguments are handled in the manner of fmt.Println.
func (l *logger) Errorln(v ...interface{}) {
	l.msg <- message{sError, fmt.Sprintln(v...)}
}

// Errorf logs with the Error severity.
// Arguments are handled in the manner of fmt.Printf.
func (l *logger) Errorf(format string, v ...interface{}) {
	l.msg <- message{sError, fmt.Sprintf(format, v...)}
}

// Fatal logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Print.
func (l *logger) Fatal(v ...interface{}) {
	l.msg <- message{sFatal, fmt.Sprint(v...)}
	os.Exit(1)
}

// Fatalln logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Println.
func (l *logger) Fatalln(v ...interface{}) {
	l.msg <- message{sFatal, fmt.Sprintln(v...)}
	os.Exit(1)
}

// Fatalf logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
func (l *logger) Fatalf(format string, v ...interface{}) {
	l.msg <- message{sFatal, fmt.Sprintf(format, v...)}
	os.Exit(1)
}
