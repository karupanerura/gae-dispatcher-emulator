package gaedispemu

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewDispatcher(t *testing.T) {
	loader := NewYAMLConfigLoader("./testdata/dispatch.yaml")
	config, err := loader.LoadConfig()
	if err != nil {
		t.Error(err)
	}

	services := map[string]*Service{
		"default": &Service{
			Name:   "default",
			Origin: mustParseURL("http://localhost:8081"),
		},
		"mobile-frontend": &Service{
			Name:   "mobile-frontend",
			Origin: mustParseURL("http://localhost:8082"),
		},
	}

	_, err = NewDispatcher(services, config)
	if err == nil {
		t.Error("Should be error because static-backend is not defined")
	}
}

func TestDispatcher(t *testing.T) {
	loader := NewYAMLConfigLoader("./testdata/dispatch.yaml")
	config, err := loader.LoadConfig()
	if err != nil {
		t.Error(err)
	}

	services := map[string]*Service{
		"default": &Service{
			Name:   "default",
			Origin: mustParseURL("http://localhost:8081"),
		},
		"mobile-frontend": &Service{
			Name:   "mobile-frontend",
			Origin: mustParseURL("http://localhost:8082"),
		},
		"static-backend": &Service{
			Name:   "static-backend",
			Origin: mustParseURL("http://localhost:8083"),
		},
	}

	dispatcher, err := NewDispatcher(services, config)
	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		Host, Path string
		Service    *Service
	}{
		{
			Host:    "simple-sample.appspot.com",
			Path:    "/",
			Service: services["default"],
		},
		{
			Host:    "simple-sample.appspot.com",
			Path:    "/register",
			Service: services["default"],
		},
		{
			Host:    "localhost",
			Path:    "/favicon.ico",
			Service: services["default"],
		},
		{
			Host:    "localhost",
			Path:    "/mobile/favicon.ico",
			Service: services["mobile-frontend"],
		},
		{
			Host:    "localhost",
			Path:    "/work/favicon.ico",
			Service: services["static-backend"],
		},
		{
			Host:    "localhost",
			Path:    "/",
			Service: nil,
		},
	}
	for _, c := range cases {
		service := dispatcher.Dispatch(c.Host, c.Path)
		if diff := cmp.Diff(c.Service, service); diff != "" {
			t.Errorf("`%s%s` is failed: diff=%s", c.Host, c.Path, diff)
		}
	}
}
