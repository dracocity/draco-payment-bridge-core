package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultListenIsUnix(t *testing.T) {
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

func TestLoad_MakesRelativeUnixPathAbsolute(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	pluginDir := filepath.Join(dir, "plugins")
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := "[[listen]]\nnetwork=\"unix\"\naddress=\"tmp/dpbc.sock\"\nplugin_dir=\"" + pluginDir + "\"\n"
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
	if !filepath.IsAbs(cfg.Listen[0].Address) {
		t.Fatalf("unix address = %q, want absolute path", cfg.Listen[0].Address)
	}
}

func TestLoad_AcceptsMultipleListenEndpoints(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	pluginDir := filepath.Join(dir, "plugins")
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := "[[listen]]\nnetwork=\"tcp\"\naddress=\":8080\"\n\n[[listen]]\nnetwork=\"unix\"\naddress=\"/tmp/app.sock\"\n\nplugin_dir=\"" + pluginDir + "\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Listen) != 2 {
		t.Fatalf("Listen length = %d, want 2", len(cfg.Listen))
	}
}

func TestLoad_RejectsInvalidListenNetwork(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	pluginDir := filepath.Join(dir, "plugins")
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := "[[listen]]\nnetwork=\"http\"\naddress=\"0.0.0.0:8080\"\nplugin_dir=\"" + pluginDir + "\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := Load(cfgPath)
	if err == nil {
		t.Fatal("Load() expected error for invalid listen network")
	}
}
