package gaedispemu

// Config is an abstruct configuration for GAE dispatch services.
type Config struct {
	Rules []ConfigRule
}

// Len is length of the config.
func (c Config) Len() int {
	return len(c.Rules)
}

// ConfigRule is an abstruct dispatch rule for GAE dispatch services.
type ConfigRule struct {
	ServiceName string
	HostPathMatcher
}
