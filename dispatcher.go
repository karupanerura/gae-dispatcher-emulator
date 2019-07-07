package gaedispemu

import "fmt"

// Dispatcher is a service dispatcher
type Dispatcher interface {
	Dispatch(host, path string) *Service
}

// NewDispatcher is a constructor of Dispatcher
func NewDispatcher(services map[string]*Service, config *Config) (Dispatcher, error) {
	for _, rule := range config.Rules {
		if _, ok := services[rule.ServiceName]; !ok {
			return nil, fmt.Errorf("Undefined backend for service: %s", rule.ServiceName)
		}
	}
	return &defaultDispatcher{services: services, config: config}, nil
}

type defaultDispatcher struct {
	services map[string]*Service
	config   *Config
}

func (d *defaultDispatcher) Dispatch(host, path string) *Service {
	for _, rule := range d.config.Rules {
		if rule.MatchHostPath(host, path) {
			service := d.services[rule.ServiceName]
			return service
		}
	}
	return nil
}
