package gaedispemu

import "testing"

func TestYAMLConfigLoader(t *testing.T) {
	loader := NewYAMLConfigLoader("./testdata/dispatch.yaml")
	config, err := loader.LoadConfig()
	if err != nil {
		t.Error(err)
	}

	testConfig(t, config)
}

func TestYAMLConfigLoaderError(t *testing.T) {
	_, err := NewYAMLConfigLoader("./testdata/naiyo-dispatch.yaml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}

	_, err = NewYAMLConfigLoader("./testdata/dispatch.xml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}

	_, err = NewYAMLConfigLoader("./testdata/invalid-dispatch.yaml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}
}
