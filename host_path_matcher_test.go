package gaedispemu

import "testing"

func TestCompileHostPathMatcher(t *testing.T) {
	t.Run("ExpectedHost", func(t *testing.T) {
		m, err := CompileHostPathMatcher("example.com/")
		if err != nil {
			t.Error(err)
		}

		if !m.MatchHostPath("example.com", "/foo") {
			t.Error("should match for example.com/foo")
		}
		if !m.MatchHostPath("example.com", "/") {
			t.Error("should match for example.com/")
		}
		if m.MatchHostPath("hoge.example.com", "/") {
			t.Error("should not match for hoge.example.com/")
		}
	})

	t.Run("MatchHost", func(t *testing.T) {
		m, err := CompileHostPathMatcher("*.example.com/")
		if err != nil {
			t.Error(err)
		}

		if !m.MatchHostPath("hoge.example.com", "/") {
			t.Error("should match for hoge.example.com/")
		}
		if !m.MatchHostPath("hoge.example.com", "/foo") {
			t.Error("should match for hoge.example.com/foo")
		}
		if m.MatchHostPath("example.com", "/") {
			t.Error("should not match for example.com/")
		}
		if m.MatchHostPath("hoge-example.com", "/") {
			t.Error("should not match for hoge.example.com/")
		}
	})

	t.Run("ExpectedPath", func(t *testing.T) {
		m, err := CompileHostPathMatcher("*/favicon.ico")
		if err != nil {
			t.Error(err)
		}

		if !m.MatchHostPath("example.com", "/favicon.ico") {
			t.Error("should match for example.com/favicon.ico")
		}
		if !m.MatchHostPath("localhost", "/favicon.ico") {
			t.Error("should match for localhost/favicon.ico")
		}
		if m.MatchHostPath("localhost", "/favicon.ico/foo") {
			t.Error("should not match for localhost/favicon.ico/foo")
		}
	})

	t.Run("MatchPath", func(t *testing.T) {
		m, err := CompileHostPathMatcher("*/service1/*")
		if err != nil {
			t.Error(err)
		}

		if !m.MatchHostPath("example.com", "/service1/") {
			t.Error("should match for example.com/service1/")
		}
		if !m.MatchHostPath("localhost", "/service1/") {
			t.Error("should match for example.com/service1/")
		}
		if !m.MatchHostPath("example.com", "/service1/foo") {
			t.Error("should match for example.com/service1/foo")
		}
		if m.MatchHostPath("example.com", "/") {
			t.Error("should not match for example.com/")
		}
		if m.MatchHostPath("example.com", "/service1") {
			t.Error("should not match for example.com/service1")
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		patterns := []string{
			"/service1/*/foo/*",
			"*/service1/*/foo/*",
			"*/service1/*/foo/",
		}
		for _, pattern := range patterns {
			_, err := CompileHostPathMatcher(pattern)
			if err == nil {
				t.Errorf("should be error for %s", pattern)
			}
		}
	})
}
