// Package filter provides port filtering logic for portwatch.
// It allows including or excluding specific ports or port ranges
// from scan results before alerting.
package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// Rule represents a single port filter rule (allow or deny).
type Rule struct {
	Low  uint16
	High uint16
}

// Filter holds compiled include/exclude rules.
type Filter struct {
	includes []Rule
	excludes []Rule
}

// New creates a Filter from include and exclude range strings.
// Each string is a comma-separated list of ports or port ranges (e.g. "80,443,8000-9000").
func New(include, exclude string) (*Filter, error) {
	f := &Filter{}
	var err error
	if include != "" {
		f.includes, err = parseRules(include)
		if err != nil {
			return nil, fmt.Errorf("include: %w", err)
		}
	}
	if exclude != "" {
		f.excludes, err = parseRules(exclude)
		if err != nil {
			return nil, fmt.Errorf("exclude: %w", err)
		}
	}
	return f, nil
}

// Allow reports whether the given port should be monitored.
func (f *Filter) Allow(port uint16) bool {
	if len(f.includes) > 0 && !matchesAny(port, f.includes) {
		return false
	}
	if matchesAny(port, f.excludes) {
		return false
	}
	return true
}

func matchesAny(port uint16, rules []Rule) bool {
	for _, r := range rules {
		if port >= r.Low && port <= r.High {
			return true
		}
	}
	return false
}

func parseRules(s string) ([]Rule, error) {
	var rules []Rule
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			sides := strings.SplitN(part, "-", 2)
			lo, err := parsePort(sides[0])
			if err != nil {
				return nil, err
			}
			hi, err := parsePort(sides[1])
			if err != nil {
				return nil, err
			}
			if lo > hi {
				return nil, fmt.Errorf("invalid range %s: low > high", part)
			}
			rules = append(rules, Rule{Low: lo, High: hi})
		} else {
			p, err := parsePort(part)
			if err != nil {
				return nil, err
			}
			rules = append(rules, Rule{Low: p, High: p})
		}
	}
	return rules, nil
}

func parsePort(s string) (uint16, error) {
	n, err := strconv.ParseUint(strings.TrimSpace(s), 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: %w", s, err)
	}
	if n == 0 {
		return 0, fmt.Errorf("port must be >= 1, got 0")
	}
	return uint16(n), nil
}
