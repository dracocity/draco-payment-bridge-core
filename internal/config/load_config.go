package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pelletier/go-toml/v2"
)

// Load loads configuration from TOML file and/or environment variables
// Priority: TOML file > environment variables > defaults
func Load(configPath string) (*Config, error) {
	cfg := &Config{}

	// Apply defaults and environment variables first (TOML overrides if present).
	if err := envconfig.Process("DPBC", cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	// Apply providers config from environment variables.
	applyProviderEnv(cfg)

	if configPath != "" {
		// Load from config file if provided path
		if data, err := os.ReadFile(configPath); err == nil {
			if err := toml.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
		}
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

// applyProviderEnv applies provider config from environment variables.
// Format: DPBC_PROVIDERS_<PROVIDER>_<KEY>=value (e.g., DPBC_PROVIDERS_NOWPAYMENTS_API_KEY=...).
// With the current call order, TOML values override env values if both exist.
func applyProviderEnv(cfg *Config) {
	const prefix = "DPBC_PROVIDERS_"
	if cfg.Providers == nil {
		cfg.Providers = make(map[string]PGConfig)
	}
	for _, env := range os.Environ() {
		key, value, ok := strings.Cut(env, "=")
		if !ok || !strings.HasPrefix(key, prefix) {
			continue
		}
		parts := strings.Split(strings.TrimPrefix(key, prefix), "_")
		if len(parts) < 2 {
			continue
		}
		provider := strings.ToLower(parts[0])
		field := strings.ToLower(strings.Join(parts[1:], "_"))
		if provider == "" || field == "" {
			continue
		}
		if cfg.Providers[provider] == nil {
			cfg.Providers[provider] = PGConfig{}
		}
		cfg.Providers[provider][field] = value
	}
}
