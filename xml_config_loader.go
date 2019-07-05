package gaedispemu

import (
	"encoding/xml"
	"os"
)

type dispatchXML struct {
	Entries []dispatchEntryXML `xml:"dispatch"`
}

type dispatchEntryXML struct {
	URL    string `xml:"url"`
	Module string `xml:"module"`
}

// XMLConfigLoader is a config loader for dispatch.xml
type XMLConfigLoader struct {
	filePath string
}

// NewXMLConfigLoader is constructor of XMLConfigLoader
func NewXMLConfigLoader(filePath string) *XMLConfigLoader {
	return &XMLConfigLoader{filePath: filePath}
}

// LoadConfig loads and parse the dispatch.xml
func (l *XMLConfigLoader) LoadConfig() (*Config, error) {
	f, err := os.Open(l.filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := xml.NewDecoder(f)

	var v dispatchXML
	err = decoder.Decode(&v)
	if err != nil {
		return nil, err
	}

	return l.transform(&v)
}

func (l *XMLConfigLoader) transform(rawConfifg *dispatchXML) (*Config, error) {
	rules := make([]ConfigRule, len(rawConfifg.Entries))
	for i, entry := range rawConfifg.Entries {
		hostPathMatcher, err := CompileHostPathMatcher(entry.URL)
		if err != nil {
			return nil, err
		}

		rules[i] = ConfigRule{
			ServiceName:     entry.Module,
			HostPathMatcher: hostPathMatcher,
		}
	}
	return &Config{Rules: rules}, nil
}
