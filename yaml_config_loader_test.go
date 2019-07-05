package gaedispemu

import "testing"

func TestYAMLConfigLoader(t *testing.T) {
	loader := NewYAMLConfigLoader("./test/dispatch.yaml")
	config, err := loader.LoadConfig()
	if err != nil {
		t.Error(err)
	}

	testConfig(t, config)
}
