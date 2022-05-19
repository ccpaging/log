// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package config

type Config struct {
	EnableConsole    bool
	ConsoleLevel     string
	ConsoleAnsiColor bool

	EnableFile      bool
	FileLevel       string
	FileLocation    string
	FileLimitSize   string
	FileBackupCount int
}

func Default() *Config {
	return &Config{
		EnableConsole:    true,
		ConsoleLevel:     "debug",
		ConsoleAnsiColor: false,
		EnableFile:       false,
		FileLevel:        "info",
		FileLocation:     "",
		FileLimitSize:    "10M",
		FileBackupCount:  7,
	}
}
