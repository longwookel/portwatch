package reporter

import (
	"fmt"
	"strings"
)

// ParseFormat converts a string to a Format, returning an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("reporter: unknown format %q (want \"text\" or \"json\")", s)
	}
}

// String returns the string representation of a Format.
func (f Format) String() string {
	return string(f)
}
