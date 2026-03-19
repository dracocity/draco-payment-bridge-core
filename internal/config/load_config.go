package config

import (
	"fmt"
	"net"
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

	var listenEntries []ListenConfig
	if strings.TrimSpace(cfg.ListenRaw) != "" {
		var err error
		listenEntries, err = parseListenRaw(cfg.ListenRaw)
		if err != nil {
			return nil, err
		}
	}

	listen, err := normalizeListenConfig(listenEntries)
	if err != nil {
		return nil, err
	}
	cfg.Listen = listen

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

func parseListenRaw(raw string) ([]ListenConfig, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("invalid listen config: value is empty")
	}

	entries := strings.Split(raw, ",")
	listen := make([]ListenConfig, 0, len(entries))
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		network, address, found := strings.Cut(entry, "@")
		if !found {
			return nil, fmt.Errorf("invalid listen entry '%s': expected network@address", entry)
		}

		network = strings.TrimSpace(network)
		address = strings.TrimSpace(address)
		if network == "" || address == "" {
			return nil, fmt.Errorf("invalid listen entry '%s': network and address are required", entry)
		}

		listen = append(listen, ListenConfig{
			Network: network,
			Address: address,
		})
	}

	if len(listen) == 0 {
		return nil, fmt.Errorf("invalid listen config: no valid entries")
	}

	return listen, nil
}

func normalizeListenConfig(listenEntries []ListenConfig) ([]ListenConfig, error) {
	if len(listenEntries) == 0 {
		listenEntries = []ListenConfig{{Network: "unix", Address: "/tmp/dpbc.sock"}}
	}

	normalized := make([]ListenConfig, 0, len(listenEntries))
	seen := make(map[string]struct{}, len(listenEntries))
	for _, entry := range listenEntries {
		network := strings.ToLower(strings.TrimSpace(entry.Network))
		address := strings.TrimSpace(entry.Address)

		switch network {
		case "tcp":
			if address == "" {
				return nil, fmt.Errorf("invalid listen config: tcp endpoint requires address")
			}
			if _, _, err := net.SplitHostPort(address); err != nil {
				return nil, fmt.Errorf("invalid tcp address '%s': must be host:port", address)
			}
		case "unix":
			if address == "" {
				return nil, fmt.Errorf("invalid listen config: unix endpoint requires address")
			}
			if !filepath.IsAbs(address) {
				if wd, err := os.Getwd(); err == nil {
					address = filepath.Join(wd, address)
				}
			}
		default:
			return nil, fmt.Errorf("invalid listen network '%s': must be tcp or unix", entry.Network)
		}

		key := network + "|" + address
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, ListenConfig{Network: network, Address: address})
	}

	if len(normalized) == 0 {
		return nil, fmt.Errorf("listen must contain at least one valid endpoint")
	}

	return normalized, nil
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
