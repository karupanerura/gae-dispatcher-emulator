package gaedispemu

import (
	"errors"
	"regexp"
)

type hostMatcher interface {
	MatchHost(host string) bool
}

var errInvalidHostPattern = errors.New("Invalid Host Pattern")

var hostPatternRegexp = regexp.MustCompile(`^\*?[^\*]+$`)

func compileHostMatcher(pattern string) (hostMatcher, error) {
	if pattern == "*" {
		return passThroughHostMatcher, nil
	}

	if !hostPatternRegexp.MatchString(pattern) {
		return nil, errInvalidHostPattern
	}

	if pattern[0] == '*' {
		suffix := pattern[1:len(pattern)]
		matcher := suffixStringMathcer{suffix: suffix}
		return &stringHostMatcher{stringMatcher: matcher}, nil
	}

	matcher := justStringMathcer{expected: pattern}
	return &stringHostMatcher{stringMatcher: matcher}, nil
}

type passThroughHostMatcherType struct{}

var passThroughHostMatcher = passThroughHostMatcherType{}

func (m passThroughHostMatcherType) MatchHost(host string) bool {
	return true
}

type stringHostMatcher struct {
	stringMatcher
}

func (m *stringHostMatcher) MatchHost(host string) bool {
	return m.Match(host)
}
