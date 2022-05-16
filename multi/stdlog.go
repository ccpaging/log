// Copyright (c) 2022-present ccpaging <ccpaging@gmail.com>. All Rights Reserved.
// See License.txt for license information.

package multi

import (
	"io"
	"log"
)

func (l *Multi) StdLog(name string) *log.Logger {
	return l.StdLogAt(Ldebug, name)
}

// StdLogAt returns *log.Logger which writes to supplied zap logger at required level.
func (l *Multi) StdLogAt(level, name string) *log.Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	n := ltoi(level)
	prefix := levelStrings[n] + name
	if ll := l.core[levelStrings[n]]; ll != nil {
		return log.New(ll.Writer(), prefix, LstdFlags)
	}

	return log.New(io.Discard, prefix, LstdFlags)
}

// NewStdLog returns a *log.Logger which writes to the supplied zap Logger at
// InfoLevel. To redirect the standard library's package-global logging
// functions, use RedirectStdLog instead.
func NewStdLog(name string) *log.Logger {
	return global.StdLog(name)
}

// NewStdLogAt returns *log.Logger which writes to supplied zap logger at
// required level.
func NewStdLogAt(level, name string) *log.Logger {
	return global.StdLogAt(level, name)
}
