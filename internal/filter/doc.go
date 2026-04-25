// Package filter implements include/exclude port filtering for portwatch.
//
// A Filter is constructed from two comma-separated rule strings: one for
// ports/ranges to include and one for ports/ranges to exclude.
//
// Rules are evaluated in order:
//  1. If include rules are specified, a port must match at least one to pass.
//  2. If exclude rules are specified, a port matching any rule is rejected.
//
// Ports are expressed as plain numbers (e.g. "80") or hyphen-separated
// inclusive ranges (e.g. "8000-9000"). Multiple rules are separated by commas.
//
// Example:
//
//	f, _ := filter.New("80,443,8000-9000", "8080")
//	f.Allow(443)  // true
//	f.Allow(8080) // false — explicitly excluded
//	f.Allow(22)   // false — not in include list
package filter
