package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestWebhookSenderPostsJSON(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := notify.NewWebhookSender(ts.URL, nil)
	e := notify.Event{Title: "port alert", Message: "443 closed", Level: notify.LevelWarn}
	if err := s.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["title"] != "port alert" {
		t.Errorf("title = %q; want %q", received["title"], "port alert")
	}
	if received["level"] != "WARN" {
		t.Errorf("level = %q; want %q", received["level"], "WARN")
	}
	if received["message"] != "443 closed" {
		t.Errorf("message = %q; want %q", received["message"], "443 closed")
	}
}

func TestWebhookSenderNon2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	s := notify.NewWebhookSender(ts.URL, nil)
	err := s.Send(notify.Event{Title: "x", Message: "y", Level: notify.LevelInfo})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestWebhookSenderNetworkError(t *testing.T) {
	s := notify.NewWebhookSender("http://127.0.0.1:1", nil)
	err := s.Send(notify.Event{Title: "x", Message: "y", Level: notify.LevelError})
	if err == nil {
		t.Fatal("expected error for unreachable endpoint")
	}
}

func TestNewWebhookSenderDefaultClient(t *testing.T) {
	s := notify.NewWebhookSender("http://example.com", nil)
	if s == nil {
		t.Fatal("expected non-nil sender")
	}
}
