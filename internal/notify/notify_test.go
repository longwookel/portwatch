package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestLevelString(t *testing.T) {
	cases := []struct {
		level notify.Level
		want  string
	}{
		{notify.LevelInfo, "INFO"},
		{notify.LevelWarn, "WARN"},
		{notify.LevelError, "ERROR"},
		{notify.Level(99), "UNKNOWN"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q; want %q", tc.level, got, tc.want)
		}
	}
}

func TestWriterSenderFormatsLine(t *testing.T) {
	var buf bytes.Buffer
	s := notify.NewWriterSender(&buf)
	e := notify.Event{Title: "port change", Message: "8080 opened", Level: notify.LevelWarn}
	if err := s.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "WARN") || !strings.Contains(got, "port change") || !strings.Contains(got, "8080 opened") {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestWriterSenderDefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil.
	s := notify.NewWriterSender(nil)
	if s == nil {
		t.Fatal("expected non-nil sender")
	}
}

func TestMultiDeliversToAll(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	s1 := notify.NewWriterSender(&buf1)
	s2 := notify.NewWriterSender(&buf2)
	m := notify.NewMulti(s1, s2)
	e := notify.Event{Title: "t", Message: "m", Level: notify.LevelInfo}
	if err := m.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf1.Len() == 0 || buf2.Len() == 0 {
		t.Error("expected both senders to receive the event")
	}
}

type errSender struct{ msg string }

func (e *errSender) Send(_ notify.Event) error { return errors.New(e.msg) }

func TestMultiJoinsErrors(t *testing.T) {
	m := notify.NewMulti(&errSender{"first"}, &errSender{"second"})
	err := m.Send(notify.Event{})
	if err == nil {
		t.Fatal("expected error from Multi")
	}
	if !strings.Contains(err.Error(), "first") || !strings.Contains(err.Error(), "second") {
		t.Errorf("error should mention both failures, got: %v", err)
	}
}

func TestMultiPartialFailureContinues(t *testing.T) {
	var buf bytes.Buffer
	good := notify.NewWriterSender(&buf)
	bad := &errSender{"oops"}
	m := notify.NewMulti(bad, good)
	_ = m.Send(notify.Event{Title: "x", Message: "y", Level: notify.LevelInfo})
	if buf.Len() == 0 {
		t.Error("good sender should still have received the event")
	}
}
