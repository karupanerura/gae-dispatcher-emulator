package gaedispemu

// ConfigLoader is an interface to load dispatch.xml or dispatch.yaml
type ConfigLoader interface {
	LoadConfig() (*Config, error)
}
