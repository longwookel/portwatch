package tagger_test

import (
	"testing"

	"github.com/user/portwatch/internal/tagger"
)

func TestTagWellKnownPort(t *testing.T) {
	tg := tagger.New(nil)
	if got := tg.Tag(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestTagHTTPS(t *testing.T) {
	tg := tagger.New(nil)
	if got := tg.Tag(443); got != "https" {
		t.Fatalf("expected https, got %q", got)
	}
}

func TestTagUnknownPort(t *testing.T) {
	tg := tagger.New(nil)
	if got := tg.Tag(1234); got != "unknown/1234" {
		t.Fatalf("expected unknown/1234, got %q", got)
	}
}

func TestCustomLabelOverridesBuiltin(t *testing.T) {
	tg := tagger.New(map[uint16]string{80: "internal-proxy"})
	if got := tg.Tag(80); got != "internal-proxy" {
		t.Fatalf("expected internal-proxy, got %q", got)
	}
}

func TestCustomLabelForUnknownPort(t *testing.T) {
	tg := tagger.New(map[uint16]string{9200: "elasticsearch"})
	if got := tg.Tag(9200); got != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %q", got)
	}
}

func TestTagAllReturnsAllPorts(t *testing.T) {
	tg := tagger.New(nil)
	ports := []uint16{22, 80, 443, 9999}
	result := tg.TagAll(ports)

	if len(result) != len(ports) {
		t.Fatalf("expected %d entries, got %d", len(ports), len(result))
	}
	if result[22] != "ssh" {
		t.Errorf("expected ssh for 22, got %q", result[22])
	}
	if result[80] != "http" {
		t.Errorf("expected http for 80, got %q", result[80])
	}
	if result[443] != "https" {
		t.Errorf("expected https for 443, got %q", result[443])
	}
	if result[9999] != "unknown/9999" {
		t.Errorf("expected unknown/9999 for 9999, got %q", result[9999])
	}
}

func TestTagAllEmptySlice(t *testing.T) {
	tg := tagger.New(nil)
	result := tg.TagAll([]uint16{})
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(result))
	}
}
