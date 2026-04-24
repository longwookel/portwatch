package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	Interval    time.Duration `yaml:"interval"`
	SnapshotDir string        `yaml:"snapshot_dir"`
	Ports       PortConfig    `yaml:"ports"`
	Alert       AlertConfig   `yaml:"alert"`
}

// PortConfig defines which ports to monitor.
type PortConfig struct {
	Protocols []string `yaml:"protocols"` // e.g. ["tcp", "udp"]
	RangeMin  uint16   `yaml:"range_min"`
	RangeMax  uint16   `yaml:"range_max"`
}

// AlertConfig defines alert output settings.
type AlertConfig struct {
	Output string `yaml:"output"` // "stdout" or a file path
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Interval:    30 * time.Second,
		SnapshotDir: "/tmp/portwatch",
		Ports: PortConfig{
			Protocols: []string{"tcp"},
			RangeMin:  1,
			RangeMax:  65535,
		},
		Alert: AlertConfig{
			Output: "stdout",
		},
	}
}

// Load reads a YAML config file from path and merges it over defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that the config values are sensible.
func (c *Config) validate() error {
	if c.Interval < time.Second {
		return ErrIntervalTooShort
	}
	if c.Ports.RangeMin > c.Ports.RangeMax {
		return ErrInvalidPortRange
	}
	return nil
}
