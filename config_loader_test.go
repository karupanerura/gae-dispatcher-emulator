package gaedispemu

import "testing"

func testConfig(t *testing.T, config *Config) {
	// test config structure
	if config == nil {
		t.Error("config should not be nil")
	}
	if config.Rules == nil {
		t.Error("config.Rules should not be nil")
	}
	if len(config.Rules) != 4 {
		t.Error("config.Rules should have 4 rules")
	}
	if config.Len() != len(config.Rules) {
		t.Error("config.Len() should be 4")
	}

	if rule := config.Rules[0]; rule.ServiceName != "default" {
		t.Errorf("config.Rules[0].ServiceName should be `default`, but got: %s", rule.ServiceName)
	} else if rule.HostPathMatcher == nil {
		t.Error("config.Rules[0].HostPathMatcher should not be nil")
	} else if !rule.HostPathMatcher.MatchHostPath("localhost", "/favicon.ico") {
		t.Error("config.Rules[0].HostPathMatcher should match for localhost/favicon.ico")
	}

	if rule := config.Rules[1]; rule.ServiceName != "default" {
		t.Errorf("config.Rules[1].ServiceName should be `default`, but got: %s", rule.ServiceName)
	} else if rule.HostPathMatcher == nil {
		t.Error("config.Rules[1].HostPathMatcher should not be nil")
	} else if !rule.HostPathMatcher.MatchHostPath("simple-sample.appspot.com", "/foo/bar") {
		t.Error("config.Rules[1].HostPathMatcher should match for simple-sample.appspot.com/foo/bar")
	}

	if rule := config.Rules[2]; rule.ServiceName != "mobile-frontend" {
		t.Errorf("config.Rules[2].ServiceName should be `mobile-frontend`, but got: %s", rule.ServiceName)
	} else if rule.HostPathMatcher == nil {
		t.Error("config.Rules[2].HostPathMatcher should not be nil")
	} else if !rule.HostPathMatcher.MatchHostPath("localhost", "/mobile/index.cgi") {
		t.Error("config.Rules[2].HostPathMatcher should match for localhost/mobile/index.cgi")
	}

	if rule := config.Rules[3]; rule.ServiceName != "static-backend" {
		t.Errorf("config.Rules[3].ServiceName should be `static-backend`, but got: %s", rule.ServiceName)
	} else if rule.HostPathMatcher == nil {
		t.Error("config.Rules[3].HostPathMatcher should not be nil")
	} else if !rule.HostPathMatcher.MatchHostPath("simple-sample.appspot.com", "/work/admin.php") {
		t.Error("config.Rules[0].HostPathMatcher should match for simple-sample.appspot.com/work/admin.php")
	}
}
