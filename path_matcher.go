package gaedispemu

import (
	"errors"
	"regexp"
)

// PathMatcher is an abstruct matcher for URL path
type pathMatcher interface {
	MatchPath(path string) bool
}

var errInvalidPathPattern = errors.New("Invalid Path Pattern")

var pathPatternRegexp = regexp.MustCompile(`^[^\*]+\*?$`)

func compilePathMatcher(pattern string) (pathMatcher, error) {
	if pattern == "" || pattern == "*" {
		return passThroughPathMatcher, nil
	}

	if !pathPatternRegexp.MatchString(pattern) {
		return nil, errInvalidPathPattern
	}

	if pattern[len(pattern)-1] == '*' {
		prefix := pattern[0 : len(pattern)-1]
		matcher := prefixStringMathcer{prefix: prefix}
		return &stringPathMatcher{stringMatcher: matcher}, nil
	}

	matcher := justStringMathcer{expected: pattern}
	return &stringPathMatcher{stringMatcher: matcher}, nil
}

type passThroughPathMatcherType struct{}

var passThroughPathMatcher = passThroughPathMatcherType{}

func (m passThroughPathMatcherType) MatchPath(path string) bool {
	return true
}

type stringPathMatcher struct {
	stringMatcher
}

func (m *stringPathMatcher) MatchPath(path string) bool {
	return m.Match(path[1:])
}
