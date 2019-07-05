package gaedispemu

import (
	"strings"
)

type stringMatcher interface {
	Match(s string) bool
}

type justStringMathcer struct {
	expected string
}

func (m justStringMathcer) Match(s string) bool {
	return s == m.expected
}

type prefixStringMathcer struct {
	prefix string
}

func (m prefixStringMathcer) Match(s string) bool {
	return strings.HasPrefix(s, m.prefix)
}

type suffixStringMathcer struct {
	suffix string
}

func (m suffixStringMathcer) Match(s string) bool {
	return strings.HasSuffix(s, m.suffix)
}
