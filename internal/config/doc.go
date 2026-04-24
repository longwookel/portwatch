// Package config provides loading and validation of portwatch configuration.
//
// Configuration is read from a YAML file and merged over built-in defaults,
// so a minimal config file only needs to specify values that differ from the
// defaults.
//
// Example usage:
//
//	cfg, err := config.Load("/etc/portwatch/config.yaml")
//	if err != nil {
//		log.Fatalf("config error: %v", err)
//	}
package config
