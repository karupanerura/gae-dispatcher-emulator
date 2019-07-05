package gaedispemu

import (
	"fmt"
	"strings"
)

// HostPathMatcher is an abstruct matcher
type HostPathMatcher interface {
	MatchHostPath(host, path string) bool
}

type genericHostPathMatcher struct {
	hostMatcher
	pathMatcher
}

func (m *genericHostPathMatcher) MatchHostPath(host, path string) bool {
	return m.MatchHost(host) && m.MatchPath(path)
}

// CompileHostPathMatcher is constructor for HostPathMatcher
func CompileHostPathMatcher(pattern string) (HostPathMatcher, error) {
	parts := strings.SplitN(pattern, "/", 2)
	host, path := parts[0], parts[1]

	hostMatcher, err := compileHostMatcher(host)
	if err != nil {
		return nil, fmt.Errorf("Invalid URL pattern: %s (%s)", pattern, err.Error())
	}

	pathMatcher, err := compilePathMatcher(path)
	if err != nil {
		return nil, fmt.Errorf("Invalid URL pattern: %s (%s)", pattern, err.Error())
	}

	return &genericHostPathMatcher{
		hostMatcher: hostMatcher,
		pathMatcher: pathMatcher,
	}, nil
}
