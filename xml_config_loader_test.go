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
