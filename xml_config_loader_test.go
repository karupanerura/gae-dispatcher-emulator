package gaedispemu

import "testing"

func TestXMLConfigLoader(t *testing.T) {
	loader := NewXMLConfigLoader("./test/dispatch.xml")
	config, err := loader.LoadConfig()
	if err != nil {
		t.Error(err)
	}

	testConfig(t, config)
}

func TestXMLConfigLoaderError(t *testing.T) {
	_, err := NewXMLConfigLoader("./test/naiyo-dispatch.xml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}

	_, err = NewXMLConfigLoader("./test/dispatch.yaml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}

	_, err = NewXMLConfigLoader("./test/invalid-dispatch.xml").LoadConfig()
	if err == nil {
		t.Error("should be error")
	}
}
