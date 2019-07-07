package gaedispemu

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDispatcher(t *testing.T) {
	loader := NewYAMLConfigLoader("./test/dispatch.yaml")
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

	dispatcher := NewDispatcher(services, config)

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

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return u
}
