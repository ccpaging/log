package multi

import (
	"bytes"
	"log"
	"testing"
)

var testFiles []string = []string{"_test.log", "_test.1.log"}

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	root := log.New(&buf, "", 0)
	l := New("main: ", root)
	l.Info("This is info")
	if want, got := "INFO main: This is info\n", buf.String(); want != got {
		t.Errorf("logger output should match %q is %q", want, got)
	}
}

func TestLoggerDebug(t *testing.T) {
	var buf bytes.Buffer
	var l Logger = New("test: ", log.New(&buf, "", log.Lshortfile|log.Lmsgprefix))
	l.Debug("This is debug")
	if want, got := "multi_test.go:24: DEBG test: This is debug\n", buf.String(); want != got {
		t.Errorf("logger debug should match %q is %q", want, got)
	}
}

func TestLoggerNew(t *testing.T) {
	l := New("new: ", log.Default())

	var buf bytes.Buffer
	ll := New("test: ", log.New(&buf, "", 0))
	dup := ll.New("temp")
	if want, got := &ll.logs, &dup.logs; want == got {
		t.Errorf("logger new should has a new core %v is %v", want, got)
	}

	p1 := &l.logs
	l.CopyFrom(dup)
	p2 := &l.logs
	if p1 != p2 {
		t.Errorf("logger copyfrom should match %v is %v", p1, p2)
	}

	l.Info("This is info")
	if want, got := "INFO new: This is info\n", buf.String(); want != got {
		t.Errorf("logger output should match %q is %q", want, got)
	}
}

func BenchmarkStdlogPrint(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := log.New(&buf, "INFO ", log.LstdFlags)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Print(testString)
	}
	b.StopTimer()
}

func BenchmarkInfo(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer

	l := New("", log.New(&buf, "", log.LstdFlags))

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Info(testString)
	}
	b.StopTimer()
}

func BenchmarkInfoPtr(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer

	l := New("", log.New(&buf, "", log.LstdFlags))
	info := l.Info

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		info(testString)
	}
	b.StopTimer()
}
