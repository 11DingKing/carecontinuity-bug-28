package serviceactivation

import (
	"fmt"
	"sync"
)

type RoutePolicy struct {
	Mode         string
	RetainFailed bool
}
type ServiceRoutes struct {
	mu     sync.RWMutex
	owners map[string]string
	policy RoutePolicy
}

func NewServiceRoutes(policy RoutePolicy) *ServiceRoutes {
	return &ServiceRoutes{owners: make(map[string]string), policy: policy}
}
func (r *ServiceRoutes) Activate(region, provider string, commit func() error) error {
	if r.policy.Mode == "eager" {
		r.mu.Lock()
		r.owners[region] = provider
		r.mu.Unlock()
		if err := commit(); err != nil {
			if r.policy.RetainFailed {
				return fmt.Errorf("activation commit: %w", err)
			}
			r.mu.Lock()
			delete(r.owners, region)
			r.mu.Unlock()
			return err
		}
		return nil
	}
	if err := commit(); err != nil {
		return fmt.Errorf("activation commit: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.owners[region] = provider
	return nil
}
func (r *ServiceRoutes) Owner(region string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.owners[region]
	return v, ok
}
