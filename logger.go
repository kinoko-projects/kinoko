/*
 * Copyright 2019 Azz. All rights reserved.
 * Use of this source code is governed by a GPL-3.0
 * license that can be found in the LICENSE file.
 */

package kinoko

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level int

const (
	_ Level = iota
	Info
	Warn
	Error
)

type Logger struct {
	name   string
	config config
}

type config struct {
	EnableColor bool
	Writer      io.Writer
}

func NewLogger(name string) *Logger {
	return &Logger{name: name, config: config{EnableColor: true, Writer: os.Stdout}}
}

func NewLoggerWithConfig(name string, enableColor bool, writer io.Writer) *Logger {
	return &Logger{name: name, config: config{EnableColor: enableColor, Writer: writer}}
}

func (l *Logger) Info(message ...interface{}) {
	l.outputMessage(Info, message...)
}

func (l *Logger) Error(message ...interface{}) {
	l.outputMessage(Error, message...)
}

func (l *Logger) Warn(message ...interface{}) {
	l.outputMessage(Warn, message...)
}

func (l *Logger) outputMessage(level Level, message ...interface{}) {
	var levelColorId int
	var levelPrefix string
	switch level {
	case Info:
		levelColorId = 36
		levelPrefix = "INFO"
	case Warn:
		levelColorId = 33
		levelPrefix = "WARN"
	case Error:
		levelColorId = 31
		levelPrefix = "ERROR"
	}
	colorChange := fmt.Sprintf("\x1b[0;%dm", levelColorId)
	colorReset := "\x1b[0m"

	var s string
	if l.config.EnableColor {
		pre := []interface{}{time.Now().Format(time.RFC3339), colorChange, levelPrefix, colorReset, l.name}
		pre = append(pre, message...)
		s = fmt.Sprintln(pre...)
	} else {
		pre := []interface{}{time.Now().Format(time.RFC3339), levelPrefix, l.name, message}
		pre = append(pre, message...)
		s = fmt.Sprintln(pre)

	}
	_, _ = l.config.Writer.Write([]byte(s))
}
