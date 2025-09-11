package billing

import "fmt"

type Registry struct {
	providers map[string]BillingProvider
}

func NewRegistry() *Registry {
	return &Registry{providers: map[string]BillingProvider{}}
}

func (r *Registry) Register(p BillingProvider) {
	r.providers[p.Key()] = p
}

func (r *Registry) Get(key string) (BillingProvider, error) {
	p, ok := r.providers[key]
	if !ok {
		return nil, fmt.Errorf("billing provider not found: %s", key)
	}
	return p, nil
}
