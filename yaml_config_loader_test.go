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

func TestYAMLConfigLoaderError(t *testing.T) {
	_, err := NewYAMLConfigLoader("./test/naiyo-dispatch.yaml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}

	_, err = NewYAMLConfigLoader("./test/dispatch.xml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}

	_, err = NewYAMLConfigLoader("./test/invalid-dispatch.yaml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}
}
