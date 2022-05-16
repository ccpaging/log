// Copyright (c) 2022-present ccpaging <ccpaging@gmail.com>. All Rights Reserved.
// See License.txt for license information.

package multi

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
	Lpanic string = "PANIC "
)

var levelStrings = []string{Ldebug, Ltrace, Linfo, Lwarn, Lerror, Lfatal}

func ltoi(s string) int {
	switch strings.ToLower(strings.Trim(s, " \r\n")) {
	case "debug", Ldebug:
		return 0
	case "trace", Ltrace:
		return 1
	case "info", Linfo:
		return 2
	case "warn", "warning", Lwarn:
		return 3
	case "err", "error", Lerror:
		return 4
	case "fatal", Lfatal:
		return 5
	default:
	}
	return 2
}

var errOutput = errors.New("No output")

type Multi struct {
	mu   sync.Mutex
	name string
	core map[string]*log.Logger
	io.Closer
}

// New creates a new Multi Level Logger. The name variable sets the
// module name will be written. The root variable sets the root logger.
// The levels enables different levels' logger.
func New(name string, root *log.Logger, levels []string) *Multi {
	if root == nil {
		root = log.Default()
		root.SetFlags(LstdFlags)
	}

	if levels == nil {
		levels = levelStrings
	}
	core := make(map[string]*log.Logger)
	for _, level := range levels {
		core[level] = root
	}
	return &Multi{
		name: name,
		core: core,
	}
}

func (l *Multi) Set(level string, logger *log.Logger) *log.Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	old := l.core[level]
	l.core[level] = logger
	return old
}

func (l *Multi) Get(level string) (logger *log.Logger, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	logger, ok = l.core[level]
	return
}

// New creates a new Multi Level Logger with a new name.
func (l *Multi) New(name string) *Multi {
	l.mu.Lock()
	defer l.mu.Unlock()

	return &Multi{
		name: name,
		core: l.core,
	}
}

// CopyFrom deletes previous loggers and copy from new logger.
func (l *Multi) CopyFrom(in *Multi) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// clear map
	for key := range l.core {
		delete(l.core, key)
	}
	// copy from
	for key, value := range in.core {
		l.core[key] = value
	}
}

// Close the logger writer if l.Closer is set.
func (l *Multi) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Closer != nil {
		l.Closer.Close()
	}

	// clear map
	for key := range l.core {
		delete(l.core, key)
	}
}

func (l *Multi) Loutput(calldepth int, level string, v ...any) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.core[level]; ok {
		return ll.Output(2+calldepth, level+l.name+fmt.Sprint(v...))
	}
	return errOutput
}

func (l *Multi) Loutputf(calldepth int, level string, format string, v ...any) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.core[level]; ok {
		return ll.Output(2+calldepth, level+l.name+fmt.Sprintf(format, v...))
	}
	return errOutput
}

func (l *Multi) Loutputln(calldepth int, level string, v ...any) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.core[level]; ok {
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

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Multi) Panic(v ...any) {
	s := fmt.Sprint(v...)
	l.Loutput(1, Lpanic, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Multi) Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	l.Loutput(1, Lpanic, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Multi) Panicln(v ...any) {
	s := fmt.Sprintln(v...)
	l.Loutput(1, Lpanic, s)
	panic(s)
}
