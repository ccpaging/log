// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ccpaging/log/file"
	"github.com/ccpaging/log/multi"
)

var levelStrings = multi.LevelStrings
var LstdFlags = log.LstdFlags

func ltoi(s string) int {
	switch strings.ToLower(strings.Trim(s, " \r\n")) {
	case "trace", multi.Ltrace:
		return 0
	case "debug", multi.Ldebug:
		return 1
	case "info", multi.Linfo:
		return 2
	case "warn", "warning", multi.Lwarn:
		return 3
	case "err", "error", multi.Lerror:
		return 4
	case "fatal", multi.Lfatal:
		return 5
	default:
	}
	return 2
}

type Builder struct {
	cw  io.Writer
	fw  *file.File
	cal int // the level index of console output
	fal int // the level index of file output
}

func NewBuilder(c *Config) *Builder {
	if c == nil {
		c = Default()
	}
	return &Builder{
		cw:  newConsoleWriter(c),
		fw:  newFileWriter(c),
		cal: ltoi(c.ConsoleLevel),
		fal: ltoi(c.FileLevel),
	}

}

func newConsoleWriter(c *Config) io.Writer {
	if !c.EnableConsole {
		return nil
	}
	if c.ConsoleAnsiColor {
		return &ansiTerm{os.Stderr}
	}
	return os.Stderr
}

func strToNumSuffix(s string, base int64) (int64, error) {
	var multi int64 = 1
	if len(s) > 1 {
		switch s[len(s)-1] {
		case 'G', 'g':
			multi *= base
			fallthrough
		case 'M', 'm':
			multi *= base
			fallthrough
		case 'K', 'k':
			multi *= base
			s = s[0 : len(s)-1]
		}
	}
	n, err := strconv.ParseInt(s, 0, 0)
	return n * multi, err
}

func newFileWriter(c *Config) *file.File {
	if !c.EnableFile {
		return nil
	}
	limitSize, _ := strToNumSuffix(c.FileLimitSize, 1024)
	if c.FileLocation == "" {
		fileName := os.Args[0]
		ext := filepath.Ext(fileName)
		c.FileLocation = fileName[0:len(fileName)-len(ext)] + "." + "log"
	}
	fw, err := file.OpenFile(c.FileLocation, limitSize, c.FileBackupCount)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Open file", err)
	}
	return fw
}

func (b *Builder) levelWriter(n int) io.Writer {
	isConsole := false
	if b.cw != nil && n >= b.cal {
		isConsole = true
	}
	isFile := false
	if b.fw != nil && n >= b.fal {
		isFile = true
	}
	if isConsole && isFile {
		return io.MultiWriter(b.cw, b.fw)
	} else if isConsole {
		return b.cw
	} else if isFile {
		return b.fw
	}
	return nil
}

func (b *Builder) Logger(name string) *multi.Multi {
	multi := multi.Omitter(name)
	for i, k := range levelStrings {
		if w := b.levelWriter(i); w != nil {
			multi.SetOutput(k, log.New(w, "", log.LstdFlags))
		}
	}
	if b.fw != nil {
		multi.Closer = b.fw
	}
	return multi
}

func (b *Builder) StdLogAt(level, name string) *log.Logger {
	n := ltoi(level)
	prefix := levelStrings[n] + name
	if w := b.levelWriter(n); w != nil {
		return log.New(w, prefix, LstdFlags)
	}

	return log.New(io.Discard, prefix, LstdFlags)
}
