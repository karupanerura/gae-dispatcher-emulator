package gaedispemu

import "net/url"

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return u
}

func mustCompileHostPathMatcher(pattern string) HostPathMatcher {
	m, err := CompileHostPathMatcher(pattern)
	if err != nil {
		panic(err)
	}

	return m
}
