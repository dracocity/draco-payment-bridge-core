package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
)

// PluginLoader loads payment bridge plugins from a directory
type PluginLoader struct {
	pluginDir string
	logger    logger.Logger

	mu      sync.RWMutex
	plugins map[string]Plugin
}

// NewPluginLoader creates a new plugin loader
func NewPluginLoader(pluginDir string) *PluginLoader {
	return &PluginLoader{
		pluginDir: pluginDir,
		logger:    logger.WithModule("plugin_loader"),
		plugins:   make(map[string]Plugin),
	}
}

// LoadPlugins loads all plugins from the plugin directory
func (l *PluginLoader) LoadPlugins() (map[string]Plugin, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.loadAllLocked()
}

// loadAllLocked scans the plugin directory and loads all .so files.
// This function is not thread-safe and must be called with a lock.
func (l *PluginLoader) loadAllLocked() (map[string]Plugin, error) {
	// Read all .so files in the plugin directory
	files, err := os.ReadDir(l.pluginDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".so" {
			continue
		}

		path := filepath.Join(l.pluginDir, file.Name())
		p, err := l.loadPlugin(path)
		if err != nil {
			l.logger.Warn("failed to load plugin", "path", path, "error", err)
			continue
		}

		l.plugins[p.GetName()] = p
		l.logger.Info("loaded plugin", "plugin_name", p.GetName(), "filename", file.Name())
	}

	return l.plugins, nil
}

// loadPlugin loads a single plugin file
func (l *PluginLoader) loadPlugin(path string) (Plugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look for the New function
	symNew, err := p.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("plugin does not export 'New' function: %w", err)
	}

	// Type assert to the factory function
	newFunc, ok := symNew.(func() Plugin)
	if !ok {
		return nil, fmt.Errorf("plugin 'New' function has wrong signature")
	}

	// Call the factory function to create the plugin
	plugin := newFunc()
	if plugin == nil {
		return nil, fmt.Errorf("plugin factory returned nil")
	}

	return plugin, nil
}
