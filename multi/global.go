// Copyright (c) 2022-present ccpaging <ccpaging@gmail.com>. All Rights Reserved.
// See License.txt for license information.

package multi

import (
	"fmt"
	"os"
)

var global *Multi = New("root ", nil, nil)

// Default returns the standard logger used by the package-level output functions.
func Default() *Multi { return global }

func InitGlobal(l *Multi) {
	global.CopyFrom(l)
}

func CloseGlobal() {
	global.Close()
	global.CopyFrom(New("root ", nil, nil))
}

func Output(calldepth int, s string) error {
	return global.Loutput(1+calldepth, Linfo, s)
}

func Trace(v ...any) {
	global.Loutput(1, Ltrace, v...)
}

func Tracef(format string, v ...any) {
	global.Loutputf(1, Ltrace, format, v...)
}

func Traceln(v ...any) {
	global.Loutputln(1, Ltrace, v...)
}

func Debug(v ...any) {
	global.Loutput(1, Ldebug, v...)
}

func Debugf(format string, v ...any) {
	global.Loutputf(1, Ldebug, format, v...)
}

func Debugln(v ...any) {
	global.Loutputln(1, Ldebug, v...)
}

func Info(v ...any) {
	global.Loutput(1, Linfo, v...)
}

func Infof(format string, v ...any) {
	global.Loutputf(1, Linfo, format, v...)
}

func Infoln(v ...any) {
	global.Loutputln(1, Linfo, v...)
}

func Warn(v ...any) {
	global.Loutput(1, Lwarn, v...)
}

func Warnf(format string, v ...any) {
	global.Loutputf(1, Lwarn, format, v...)
}

func Warnln(v ...any) {
	global.Loutputln(1, Lwarn, v...)
}

func Error(v ...any) {
	global.Loutput(1, Lerror, v...)
}

func Errorf(format string, v ...any) {
	global.Loutputf(1, Lerror, format, v...)
}

func Errorln(v ...any) {
	global.Loutputln(1, Lerror, v...)
}

func Print(v ...any) {
	global.Loutput(1, Linfo, v...)
}

func Printf(format string, v ...any) {
	global.Loutputf(1, Linfo, format, v...)
}

func Println(v ...any) {
	global.Loutputln(1, Linfo, v...)
}

func Fatal(v ...any) {
	global.Loutput(1, Lfatal, v...)
	global.Close()
	os.Exit(1)
}

func Fatalf(format string, v ...any) {
	global.Loutputf(1, Lfatal, format, v...)
	global.Close()
	os.Exit(1)
}

func Fatalln(v ...any) {
	global.Loutputln(1, Lfatal, v...)
	global.Close()
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...any) {
	s := fmt.Sprint(v...)
	global.Loutput(1, Lpanic, s)
	global.Close()
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	global.Loutput(1, Lpanic, s)
	global.Close()
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...any) {
	s := fmt.Sprintln(v...)
	global.Loutputf(1, Lpanic, s)
	global.Close()
	panic(s)
}
