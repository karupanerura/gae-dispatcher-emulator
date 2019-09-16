package gaedispemu

import "testing"

func TestXMLConfigLoader(t *testing.T) {
	loader := NewXMLConfigLoader("./testdata/dispatch.xml")
	config, err := loader.LoadConfig()
	if err != nil {
		t.Error(err)
	}

	testConfig(t, config)
}

func TestXMLConfigLoaderError(t *testing.T) {
	_, err := NewXMLConfigLoader("./testdata/naiyo-dispatch.xml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}

	_, err = NewXMLConfigLoader("./testdata/dispatch.yaml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}

	_, err = NewXMLConfigLoader("./testdata/invalid-dispatch.xml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}
}
