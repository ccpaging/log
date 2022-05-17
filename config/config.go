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

type Config struct {
	EnableConsole    bool
	ConsoleLevel     string
	ConsoleAnsiColor bool

	EnableFile      bool
	FileLevel       string
	FileLocation    string
	FileLimitSize   string
	FileBackupCount int

	isCalc bool
	cw     io.Writer
	fw     *file.File
	cal    int // the level index of console output
	fal    int // the level index of file output
}

func Default() *Config {
	c := &Config{
		EnableConsole:    true,
		ConsoleLevel:     multi.Ldebug,
		ConsoleAnsiColor: false,
		EnableFile:       false,
		FileLevel:        multi.Linfo,
		FileLocation:     "",
		FileLimitSize:    "10M",
		FileBackupCount:  7,
	}
	return c.calc()
}

type Logger struct {
	*multi.Multi
}

var levelStrings = multi.LevelStrings

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

func levelWriter(n int, c *Config) io.Writer {
	isConsole := false
	if c.cw != nil && n >= c.cal {
		isConsole = true
	}
	isFile := false
	if c.fw != nil && n >= c.fal {
		isFile = true
	}
	if isConsole && isFile {
		return io.MultiWriter(c.cw, c.fw)
	} else if isConsole {
		return c.cw
	} else if isFile {
		return c.fw
	}
	return nil
}

func (c *Config) Close() error {
	if c.fw != nil {
		c.fw.Close()
	}
	// Clear data
	c.fw = nil
	return nil
}

func (c *Config) calc() *Config {
	if c.isCalc {
		return c
	}
	c.cw = newConsoleWriter(c)
	c.fw = newFileWriter(c)
	c.cal = ltoi(c.ConsoleLevel)
	c.fal = ltoi(c.FileLevel)
	return c
}

func (c *Config) MakeLogger(name string) *multi.Multi {
	c.calc()

	multi := multi.Omitter(name)
	for i, k := range levelStrings {
		if w := levelWriter(i, c); w != nil {
			multi.SetOutput(k, log.New(w, "", log.LstdFlags))
		}
	}
	multi.Closer = c
	return multi
}

func (c *Config) StdLog(name string) *log.Logger {
	return c.StdLogAt(multi.Ldebug, name)
}

// StdLogAt returns *log.Logger which writes to supplied zap logger at required level.
func (c *Config) StdLogAt(level, name string) *log.Logger {
	c.calc()

	n := ltoi(level)
	if w := levelWriter(n, c); w != nil {
		return log.New(w, levelStrings[n]+name, log.LstdFlags)
	}

	return log.New(io.Discard, levelStrings[n]+name, log.LstdFlags)
}
