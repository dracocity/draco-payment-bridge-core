package pg

import "sync"

// Registry holds loaded provider payment gateways.
type Registry struct {
	mu  sync.RWMutex
	pgs map[string]PaymentGateway
}

func NewRegistry() *Registry {
	return &Registry{pgs: make(map[string]PaymentGateway)}
}

func (r *Registry) Register(p PaymentGateway) {
	if p == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pgs[p.Name()] = p
}

func (r *Registry) Get(name string) (PaymentGateway, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.pgs[name]
	return b, ok
}

func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.pgs))
	for name := range r.pgs {
		names = append(names, name)
	}
	return names
}
