package config

import "errors"

// Sentinel errors for config validation.
var (
	// ErrIntervalTooShort is returned when the scan interval is less than 1 second.
	ErrIntervalTooShort = errors.New("config: interval must be at least 1 second")

	// ErrInvalidPortRange is returned when range_min > range_max.
	ErrInvalidPortRange = errors.New("config: port range_min must be <= range_max")
)
