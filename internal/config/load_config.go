package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/pelletier/go-toml/v2"
)

// Load loads configuration from TOML file and/or environment variables
// Priority: environment variables > TOML file > defaults
func Load(configPath string) (*Config, error) {
	cfg := &Config{}

	if configPath != "" {
		// Load from config file if provided path
		if data, err := os.ReadFile(configPath); err == nil {
			if err := toml.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
		}
	}

	// Environment variables have highest priority, override TOML and apply defaults
	if err := envconfig.Process("DPBC", cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	// Expand relative paths to absolute paths
	if !filepath.IsAbs(cfg.PluginDir) {
		if wd, err := os.Getwd(); err == nil {
			cfg.PluginDir = filepath.Join(wd, cfg.PluginDir)
		}
	}

	// Validate plugin directory
	if dirinfo, err := os.Stat(cfg.PluginDir); err != nil || !dirinfo.IsDir() {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(cfg.PluginDir, 0o755); err != nil {
				return nil, fmt.Errorf("failed to create plugin directory '%s': %w", cfg.PluginDir, err)
			}
		} else if err == nil {
			return nil, fmt.Errorf("plugin path '%s' is not a directory", cfg.PluginDir)
		} else {
			return nil, fmt.Errorf("plugin directory error '%s': %w", cfg.PluginDir, err)
		}
	}

	return cfg, nil
}
