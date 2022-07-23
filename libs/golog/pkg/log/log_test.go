package log

import (
	"bytes"
	"errors"

	"strings"
	"testing"
)

func TestTrace(t *testing.T) {
	commonTest(t, "trace")
}

func TestDebug(t *testing.T) {
	commonTest(t, "debug")
}

func TestInfo(t *testing.T) {
	commonTest(t, "info")
}

func TestWarn(t *testing.T) {
	commonTest(t, "warn")
}

func commonTest(t *testing.T, levelToTest string) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		l, err := New(Trace, out)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		switch levelToTest {
		case "trace":
			l.Trace("")
		case "debug":
			l.Debug("")
		case "info":
			l.Info("")
		case "warn":
			l.Warn("")
		}
		got := out.String()
		c := strings.Contains(got, `{"level":"`+levelToTest+`"}`)
		if !c {
			t.Errorf("invalid log output, expected level info, but we got:  %s", got)
		}
	})
	t.Run("with-fields", func(t *testing.T) {
		out := &bytes.Buffer{}
		l, err := New(Trace, out)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		switch levelToTest {
		case "trace":
			l.Trace("message", LogParam{Key: "key1", Value: "str"}, LogParam{Key: "key2", Value: 2})
		case "debug":
			l.Debug("message", LogParam{Key: "key1", Value: "str"}, LogParam{Key: "key2", Value: 2})
		case "info":
			l.Info("message", LogParam{Key: "key1", Value: "str"}, LogParam{Key: "key2", Value: 2})
		case "warn":
			l.Warn("message", LogParam{Key: "key1", Value: "str"}, LogParam{Key: "key2", Value: 2})
		}
		got := out.String()
		c := strings.Contains(got, `{"level":"`+levelToTest+`","key1":"str","key2":2,"message":"message"}`)
		if !c {
			t.Errorf("invalid log output, expected level info, but we got:  %s", got)
		}
	})
}

func TestError(t *testing.T) {
	commonTestWithStack(t, "error")
}

func commonTestWithStack(t *testing.T, levelToTest string) {
	out := &bytes.Buffer{}
	l, _ := New(Trace, out)
	err := errors.New("error message")
	switch levelToTest {
	case "error":
		l.Error(err, "")
	case "fatal":
		l.Fatal(err, "")
	case "panic":
		l.Panic(err, "")
	}
	got := out.String()
	c := strings.Contains(got, "stack")
	if !c {
		t.Errorf("invalid log output, expected stack trace, but we got:  %s", got)
	} else {
		t.Logf("got stack trace: %s", got)
	}
}
