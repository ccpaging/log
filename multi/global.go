// Copyright (c) 2022-present ccpaging <ccpaging@gmail.com>. All Rights Reserved.
// See License.txt for license information.

package multi

import (
	"os"
)

var global *Multi = Default()

func Global() *Multi { return global }

func Redirect(l *Multi) {
	global.CopyFrom(l)
}

func Restore() {
	global.Close()
	global.CopyFrom(Default())
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
