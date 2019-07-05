package gaedispemu

// Config is an abstruct configuration for GAE dispatch services.
type Config struct {
	Rules []ConfigRule
}

// ConfigRule is an abstruct dispatch rule for GAE dispatch services.
type ConfigRule struct {
	ServiceName string
	HostPathMatcher
}
