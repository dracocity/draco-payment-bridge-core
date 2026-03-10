package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_ListenFromEnvWhenTomlListenOmitted(t *testing.T) {
	pluginDir := filepath.Join(t.TempDir(), "plugins")
	t.Setenv("DPBC_PLUGIN_DIR", pluginDir)
	t.Setenv("DPBC_LISTEN", "tcp@127.0.0.1:8080,unix@/tmp/dpbc.sock")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Listen) != 2 {
		t.Fatalf("Listen length = %d, want 2", len(cfg.Listen))
	}

	if cfg.Listen[0].Network != "tcp" || cfg.Listen[0].Address != "127.0.0.1:8080" {
		t.Fatalf("Listen[0] = %#v, want tcp@127.0.0.1:8080", cfg.Listen[0])
	}

	if cfg.Listen[1].Network != "unix" || cfg.Listen[1].Address != "/tmp/dpbc.sock" {
		t.Fatalf("Listen[1] = %#v, want unix@/tmp/dpbc.sock", cfg.Listen[1])
	}
}

func TestLoad_TomlListenOverridesEnvListen(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "plugins")
	t.Setenv("DPBC_LISTEN", "tcp@127.0.0.1:8080")

	cfgPath := filepath.Join(dir, "config.toml")
	content := "plugin_dir=\"" + pluginDir + "\"\n[[listen]]\nnetwork=\"unix\"\naddress=\"/tmp/from-toml.sock\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Listen) != 1 {
		t.Fatalf("Listen length = %d, want 1", len(cfg.Listen))
	}

	if cfg.Listen[0].Network != "unix" || cfg.Listen[0].Address != "/tmp/from-toml.sock" {
		t.Fatalf("Listen[0] = %#v, want unix@/tmp/from-toml.sock", cfg.Listen[0])
	}
}

func TestLoad_ReturnsErrorForInvalidListenEnvFormat(t *testing.T) {
	pluginDir := filepath.Join(t.TempDir(), "plugins")
	t.Setenv("DPBC_PLUGIN_DIR", pluginDir)
	t.Setenv("DPBC_LISTEN", "tcp127.0.0.1:8080")

	_, err := Load("")
	if err == nil {
		t.Fatal("Load() expected error for invalid DPBC_LISTEN format")
	}

	if !strings.Contains(err.Error(), "expected network@address") {
		t.Fatalf("error = %v, want message containing expected network@address", err)
	}
}

func TestLoad_DefaultListenIsUnixWhenNoTomlOrEnv(t *testing.T) {
	pluginDir := filepath.Join(t.TempDir(), "plugins")
	t.Setenv("DPBC_PLUGIN_DIR", pluginDir)
	t.Setenv("DPBC_LISTEN", "")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Listen) != 1 {
		t.Fatalf("Listen length = %d, want 1", len(cfg.Listen))
	}

	if cfg.Listen[0].Network != "unix" || cfg.Listen[0].Address != "/tmp/dpbc.sock" {
		t.Fatalf("Listen = %#v, want unix:/tmp/dpbc.sock", cfg.Listen)
	}
}

func TestLoad_MakesRelativeUnixPathAbsoluteFromEnv(t *testing.T) {
	pluginDir := filepath.Join(t.TempDir(), "plugins")
	t.Setenv("DPBC_PLUGIN_DIR", pluginDir)
	t.Setenv("DPBC_LISTEN", "unix@tmp/dpbc.sock")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Listen) != 1 {
		t.Fatalf("Listen length = %d, want 1", len(cfg.Listen))
	}

	if !filepath.IsAbs(cfg.Listen[0].Address) {
		t.Fatalf("unix address = %q, want absolute path", cfg.Listen[0].Address)
	}
}
