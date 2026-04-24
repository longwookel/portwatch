package reporter_test

import (
	"testing"

	"github.com/user/portwatch/internal/reporter"
)

func TestParseFormatText(t *testing.T) {
	for _, input := range []string{"text", "TEXT", "Text", ""} {
		f, err := reporter.ParseFormat(input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", input, err)
		}
		if f != reporter.FormatText {
			t.Errorf("ParseFormat(%q) = %v, want FormatText", input, f)
		}
	}
}

func TestParseFormatJSON(t *testing.T) {
	for _, input := range []string{"json", "JSON", "Json"} {
		f, err := reporter.ParseFormat(input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", input, err)
		}
		if f != reporter.FormatJSON {
			t.Errorf("ParseFormat(%q) = %v, want FormatJSON", input, f)
		}
	}
}

func TestParseFormatUnknown(t *testing.T) {
	_, err := reporter.ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestFormatString(t *testing.T) {
	if reporter.FormatText.String() != "text" {
		t.Errorf("FormatText.String() = %q, want \"text\"", reporter.FormatText.String())
	}
	if reporter.FormatJSON.String() != "json" {
		t.Errorf("FormatJSON.String() = %q, want \"json\"", reporter.FormatJSON.String())
	}
}
