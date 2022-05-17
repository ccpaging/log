// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package syslog_test

import (
	"fmt"
	"log"

	"github.com/ccpaging/log/syslog"
)

func ExampleDial() {
	l, err := syslog.New("tcp", "localhost:1234",
		syslog.LOG_WARNING|syslog.LOG_DAEMON, "demotag")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(l.Writer(), "This is a daemon warning with demotag.")
	l.Emerg("And this is a daemon emergency with demotag.")
}