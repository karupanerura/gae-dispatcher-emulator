package gaedispemu

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type dispatchYAML struct {
	Entries []dispatchEntryYAML `yaml:"dispatch"`
}

type dispatchEntryYAML struct {
	URL         string `yaml:"url"`
	ServiceName string `yaml:"service"`
}

// YAMLConfigLoader is a config loader for dispatch.yaml
type YAMLConfigLoader struct {
	filePath string
}

// NewYAMLConfigLoader is constructor of YAMLConfigLoader
func NewYAMLConfigLoader(filePath string) *YAMLConfigLoader {
	return &YAMLConfigLoader{filePath: filePath}
}

// LoadConfig loads and parse the dispatch.yaml
func (l *YAMLConfigLoader) LoadConfig() (*Config, error) {
	f, err := os.Open(l.filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)

	var v dispatchYAML
	err = decoder.Decode(&v)
	if err != nil {
		return nil, err
	}

	return l.transform(&v)
}

func (l *YAMLConfigLoader) transform(rawConfifg *dispatchYAML) (*Config, error) {
	rules := make([]ConfigRule, len(rawConfifg.Entries))
	for i, entry := range rawConfifg.Entries {
		hostPathMatcher, err := CompileHostPathMatcher(entry.URL)
		if err != nil {
			return nil, err
		}

		rules[i] = ConfigRule{
			ServiceName:     entry.ServiceName,
			HostPathMatcher: hostPathMatcher,
		}
	}
	return &Config{Rules: rules}, nil
}
