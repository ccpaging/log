// Copyright (c) 2022-present ccpaging <ccpaging@gmail.com>. All Rights Reserved.
// See License.txt for license information.

package multi

type Logger interface {
	// Error is equivalent to Print() and logs the message at level Error.
	Error(v ...any)
	// Errorf is equivalent to Printf() and logs the message at level Error.
	Errorf(format string, v ...any)
	// Errorln is equivalent to Println() and logs the message at level Error.
	Errorln(v ...any)

	// Warn is equivalent to Print() and logs the message at level Warning.
	Warn(v ...any)
	// Warnf is equivalent to Printf() and logs the message at level Warning.
	Warnf(format string, v ...any)
	// Warnln is equivalent to Println() and logs the message at level Warning.
	Warnln(v ...any)

	// Info is equivalent to Print() and logs the message at level Info.
	Info(v ...any)
	// Infof is equivalent to Printf() and logs the message at level Info.
	Infof(format string, v ...any)
	// Infoln is equivalent to Println() and logs the message at level Info.
	Infoln(v ...any)

	// Debug is equivalent to Print() and logs the message at level Debug.
	Debug(v ...any)
	// Debugf is equivalent to Printf() and logs the message at level Debug.
	Debugf(format string, v ...any)
	// Debugln is equivalent to Println() and logs the message at level Debug.
	Debugln(v ...any)

	// Trace is equivalent to Print() and logs the message at level Trace.
	Trace(v ...any)
	// Tracef is equivalent to Printf() and logs the message at level Trace.
	Tracef(format string, v ...any)
	// Traceln is equivalent to Println() and logs the message at level Trace.
	Traceln(v ...any)
}
