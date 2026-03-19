package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_DefaultListenWhenUnset(t *testing.T) {
	pluginDir := filepath.Join(t.TempDir(), "plugins")
	t.Setenv("DPBC_PLUGIN_DIR", pluginDir)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Listen) != 1 {
		t.Fatalf("len(cfg.Listen) = %d, want 1", len(cfg.Listen))
	}

	if cfg.Listen[0].Network != "unix" || cfg.Listen[0].Address != "/tmp/dpbc.sock" {
		t.Fatalf("cfg.Listen[0] = %#v, want unix@/tmp/dpbc.sock", cfg.Listen[0])
	}

	if _, err := os.Stat(pluginDir); err != nil {
		t.Fatalf("plugin directory should be created, stat error = %v", err)
	}
}

func TestLoad_ListenFromTomlRaw(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "plugins")
	cfgPath := filepath.Join(dir, "config.toml")
	content := "plugin_dir=\"" + pluginDir + "\"\nlisten=\"tcp@127.0.0.1:8080,unix@/tmp/from-toml.sock\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Listen) != 2 {
		t.Fatalf("len(cfg.Listen) = %d, want 2", len(cfg.Listen))
	}

	if cfg.Listen[0].Network != "tcp" || cfg.Listen[0].Address != "127.0.0.1:8080" {
		t.Fatalf("cfg.Listen[0] = %#v, want tcp@127.0.0.1:8080", cfg.Listen[0])
	}

	if cfg.Listen[1].Network != "unix" || cfg.Listen[1].Address != "/tmp/from-toml.sock" {
		t.Fatalf("cfg.Listen[1] = %#v, want unix@/tmp/from-toml.sock", cfg.Listen[1])
	}
}

func TestLoad_MakesRelativeUnixPathAbsoluteFromTomlListen(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "plugins")
	cfgPath := filepath.Join(dir, "config.toml")
	content := "plugin_dir=\"" + pluginDir + "\"\nlisten=\"unix@tmp/dpbc.sock\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Listen) != 1 {
		t.Fatalf("len(cfg.Listen) = %d, want 1", len(cfg.Listen))
	}

	if !filepath.IsAbs(cfg.Listen[0].Address) {
		t.Fatalf("cfg.Listen[0].Address = %q, want absolute path", cfg.Listen[0].Address)
	}
}

func TestLoad_DeduplicatesListenEntries(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "plugins")
	cfgPath := filepath.Join(dir, "config.toml")
	content := "plugin_dir=\"" + pluginDir + "\"\nlisten=\"tcp@127.0.0.1:8080,tcp@127.0.0.1:8080,unix@/tmp/dpbc.sock,unix@/tmp/dpbc.sock\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Listen) != 2 {
		t.Fatalf("len(cfg.Listen) = %d, want 2", len(cfg.Listen))
	}

	if cfg.Listen[0].Network != "tcp" || cfg.Listen[0].Address != "127.0.0.1:8080" {
		t.Fatalf("cfg.Listen[0] = %#v, want tcp@127.0.0.1:8080", cfg.Listen[0])
	}

	if cfg.Listen[1].Network != "unix" || cfg.Listen[1].Address != "/tmp/dpbc.sock" {
		t.Fatalf("cfg.Listen[1] = %#v, want unix@/tmp/dpbc.sock", cfg.Listen[1])
	}
}

func TestLoad_ReturnsErrorForInvalidListenEntryFormat(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "plugins")
	cfgPath := filepath.Join(dir, "config.toml")
	content := "plugin_dir=\"" + pluginDir + "\"\nlisten=\"tcp127.0.0.1:8080\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := Load(cfgPath)
	if err == nil {
		t.Fatal("Load() expected error for invalid listen format")
	}

	if !strings.Contains(err.Error(), "expected network@address") {
		t.Fatalf("error = %v, want message containing expected network@address", err)
	}
}

func TestLoad_ReturnsErrorForInvalidTCPAddress(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "plugins")
	cfgPath := filepath.Join(dir, "config.toml")
	content := "plugin_dir=\"" + pluginDir + "\"\nlisten=\"tcp@127.0.0.1\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := Load(cfgPath)
	if err == nil {
		t.Fatal("Load() expected error for invalid tcp address")
	}

	if !strings.Contains(err.Error(), "must be host:port") {
		t.Fatalf("error = %v, want message containing must be host:port", err)
	}
}

func TestLoad_AppliesProviderEnv(t *testing.T) {
	pluginDir := filepath.Join(t.TempDir(), "plugins")
	t.Setenv("DPBC_PLUGIN_DIR", pluginDir)
	t.Setenv("DPBC_PROVIDERS_NOWPAYMENTS_API_KEY", "env-key")
	t.Setenv("DPBC_PROVIDERS_NOWPAYMENTS_IPN_SECRET", "env-secret")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Providers["nowpayments"]["api_key"] != "env-key" {
		t.Fatalf("api_key = %q, want env-key", cfg.Providers["nowpayments"]["api_key"])
	}

	if cfg.Providers["nowpayments"]["ipn_secret"] != "env-secret" {
		t.Fatalf("ipn_secret = %q, want env-secret", cfg.Providers["nowpayments"]["ipn_secret"])
	}
}

func TestLoad_TomlProvidersOverrideProviderEnv(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "plugins")
	t.Setenv("DPBC_PROVIDERS_NOWPAYMENTS_API_KEY", "env-key")
	cfgPath := filepath.Join(dir, "config.toml")
	content := "plugin_dir=\"" + pluginDir + "\"\n[providers.nowpayments]\napi_key=\"toml-key\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Providers["nowpayments"]["api_key"] != "toml-key" {
		t.Fatalf("api_key = %q, want toml-key", cfg.Providers["nowpayments"]["api_key"])
	}
}
