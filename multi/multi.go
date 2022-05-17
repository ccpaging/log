// Copyright (c) 2022-present ccpaging <ccpaging@gmail.com>. All Rights Reserved.
// See License.txt for license information.

package multi

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var (
	LstdFlags = log.LstdFlags | log.Lmsgprefix
)

const (
	Ltrace string = "TRAC "
	Ldebug string = "DEBG "
	Linfo  string = "INFO "
	Lwarn  string = "WARN "
	Lerror string = "ERROR "
	Lfatal string = "FATAL "
)

var LevelStrings = []string{Ltrace, Ldebug, Linfo, Lwarn, Lerror, Lfatal}

var errOutput = errors.New("No output")

type Multi struct {
	mu   sync.Mutex
	name string
	logs map[string]Outputter
	io.Closer
}

// New creates a new Multi Level Logger. The name variable sets the
// module name will be written. The root variable sets the root logger.
// The levels enables different levels' logger.
func New(name string, output Outputter) *Multi {
	logs := make(map[string]Outputter)
	for _, level := range LevelStrings {
		logs[level] = output
	}
	return &Multi{
		name: name,
		logs: logs,
	}
}

// Default returns the standard logger used by the package-level output functions.
func Default() *Multi {
	return New("root ", log.Default())
}

func Omitter(name string) *Multi {
	return &Multi{
		name: name,
		logs: make(map[string]Outputter),
	}
}

// New creates a new Multi Level Logger with a new name.
func (l *Multi) SetOutput(level string, out Outputter) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs[level] = out
}

// New creates a new Multi Level Logger with a new name.
func (l *Multi) New(name string) *Multi {
	l.mu.Lock()
	defer l.mu.Unlock()

	return &Multi{
		name: name,
		logs: l.logs,
	}
}

// CopyFrom deletes previous loggers and copy from new logger.
func (l *Multi) CopyFrom(in *Multi) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// clear map
	for key := range l.logs {
		delete(l.logs, key)
	}
	// copy from
	for key, value := range in.logs {
		l.logs[key] = value
	}
	l.Closer = in.Closer
}

// Close the logger writer if l.Clogser is set.
func (l *Multi) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Closer != nil {
		l.Closer.Close()
	}

	// clear map
	for key := range l.logs {
		delete(l.logs, key)
	}
}

func (l *Multi) Loutput(calldepth int, level string, v ...any) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.logs[level]; ok {
		return ll.Output(2+calldepth, level+l.name+fmt.Sprint(v...))
	}
	return errOutput
}

func (l *Multi) Loutputf(calldepth int, level string, format string, v ...any) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.logs[level]; ok {
		return ll.Output(2+calldepth, level+l.name+fmt.Sprintf(format, v...))
	}
	return errOutput
}

func (l *Multi) Loutputln(calldepth int, level string, v ...any) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.logs[level]; ok {
		return ll.Output(2+calldepth, level+l.name+fmt.Sprintln(v...))
	}
	return errOutput
}

func (l *Multi) Output(calldepth int, s string) error {
	return l.Loutput(1+calldepth, Linfo, s)
}

func (l *Multi) Trace(v ...any) {
	l.Loutput(1, Ltrace, v...)
}

func (l *Multi) Tracef(format string, v ...any) {
	l.Loutputf(1, Ltrace, format, v...)
}

func (l *Multi) Traceln(v ...any) {
	l.Loutputln(1, Ltrace, v...)
}

func (l *Multi) Debug(v ...any) {
	l.Loutput(1, Ldebug, v...)
}

func (l *Multi) Debugf(format string, v ...any) {
	l.Loutputf(1, Ldebug, format, v...)
}

func (l *Multi) Debugln(v ...any) {
	l.Loutputln(1, Ldebug, v...)
}

func (l *Multi) Info(v ...any) {
	l.Loutput(1, Linfo, v...)
}

func (l *Multi) Infof(format string, v ...any) {
	l.Loutputf(1, Linfo, format, v...)
}

func (l *Multi) Infoln(v ...any) {
	l.Loutputln(1, Linfo, v...)
}

func (l *Multi) Warn(v ...any) {
	l.Loutput(1, Lwarn, v...)
}

func (l *Multi) Warnf(format string, v ...any) {
	l.Loutputf(1, Lwarn, format, v...)
}

func (l *Multi) Warnln(v ...any) {
	l.Loutputln(1, Lwarn, v...)
}

func (l *Multi) Error(v ...any) {
	l.Loutput(1, Lerror, v...)
}

func (l *Multi) Errorf(format string, v ...any) {
	l.Loutputf(1, Lerror, format, v...)
}

func (l *Multi) Errorln(v ...any) {
	l.Loutputln(1, Lerror, v...)
}

func (l *Multi) Print(v ...any) {
	l.Loutput(1, Linfo, v...)
}

func (l *Multi) Printf(format string, v ...any) {
	l.Loutputf(1, Linfo, format, v...)
}

func (l *Multi) Println(v ...any) {
	l.Loutputln(1, Linfo, v...)
}

func (l *Multi) Fatal(v ...any) {
	l.Loutput(1, Lfatal, v...)
	l.Close()
	os.Exit(1)
}

func (l *Multi) Fatalf(format string, v ...any) {
	l.Loutputf(1, Lfatal, format, v...)
	l.Close()
	os.Exit(1)
}

func (l *Multi) Fatalln(v ...any) {
	l.Loutputln(1, Lfatal, v...)
	l.Close()
	os.Exit(1)
}
