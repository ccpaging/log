package config_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/ccpaging/log/config"
	"github.com/ccpaging/log/multi"
)

var moduleLog = multi.Global().WithName("[module] ")

const testLogFile = "_test.log"

var builder = config.NewBuilder(&config.Config{
	EnableConsole:    true,
	ConsoleLevel:     "debug",
	ConsoleAnsiColor: true,

	EnableFile:      true,
	FileLevel:       "info",
	FileLocation:    testLogFile,
	FileLimitSize:   "1024k",
	FileBackupCount: 2,
})

func init() {
	os.Remove(testLogFile)
}

func removeFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	if err != nil && t != nil {
		t.Errorf("remove (%q): %s", filename, err)
	}
}

func TestConfig(t *testing.T) {
	moduleLog.Debug("debug log. ", "key=", "value")
	moduleLog.Info("info log. ", "key=", "value")
	moduleLog.Warn("warning log. ", "key=", "value")
	moduleLog.Error("error log. ", "key=", "value")

	logger := builder.Logger("test: ")

	defer removeFile(t, testLogFile)

	logger.Trace("trace log. ", "This should not be displayed")
	logger.Debug("debug log. ", "key=", "value")
	logger.Info("info log. ", "key=", "value")
	logger.Warn("warning log. ", "key=", "value")
	logger.Error("error log. ", "key=", "value")

	multi.Redirect(logger)

	multi.Debug("debug log. ", "key=", "value")
	multi.Info("info log. ", "key=", "value")
	multi.Warn("warning log. ", "key=", "value")
	multi.Error("error log. ", "key=", "value")

	moduleLog.Debug("debug log. ", "key=", "value")
	moduleLog.Info("info log. ", "key=", "value")
	moduleLog.Warn("warning log. ", "key=", "value")
	moduleLog.Error("error log. ", "key=", "value")

	multi.Restore()
	moduleLog.Error("error log. ", "key=", "value")
}

func TestStdLogAt(t *testing.T) {
	var buf bytes.Buffer

	logAtInfo := builder.StdLogAt("info", "test: ")
	logAtInfo.SetFlags(log.Lshortfile)
	logAtInfo.SetOutput(&buf)
	logAtInfo.Println("This is stdlog's Println")
	if want, got := "INFO test: config_test.go:79: This is stdlog's Println\n", buf.String(); want != got {
		t.Errorf("\nwant: %q\ngot:  %q", want, got)
	}
	buf.Reset()

	logAtDebug := builder.StdLogAt("trace", "test: ")
	logAtDebug.Println("This is stdlog's trace")
	if logAtDebug.Writer() != io.Discard {
		t.Errorf("\nwant empty\ngot: %q", buf.String())
	}
	buf.Reset()
}
