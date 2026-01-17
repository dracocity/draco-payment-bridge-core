package bridge

import "sync"

// Registry holds loaded provider bridges.
type Registry struct {
	mu      sync.RWMutex
	bridges map[string]Bridge
}

func NewRegistry() *Registry {
	return &Registry{bridges: make(map[string]Bridge)}
}

func (r *Registry) Register(b Bridge) {
	if b == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bridges[b.Name()] = b
}

func (r *Registry) Get(name string) (Bridge, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.bridges[name]
	return b, ok
}

func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.bridges))
	for name := range r.bridges {
		names = append(names, name)
	}
	return names
}
