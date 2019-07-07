package gaedispemu

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/google/go-cmp/cmp"
)

func TestProxyHandlerWithReporter(t *testing.T) {
	reporter := ErrorReporterFunc(func(err error) {})
	handler := NewProxyHandlerWithReporter(nil, reporter).(*proxyHandler)
	if reflect.ValueOf(handler.errorReporter).Pointer() != reflect.ValueOf(reporter).Pointer() {
		t.Errorf("should set expected reporter")
	}
}

func TestProxyHandler(t *testing.T) {
	defaultBackend := httptest.NewServer(getBackendHandler("default"))
	defer defaultBackend.Close()

	fooBackend := httptest.NewServer(getBackendHandler("foo"))
	defer fooBackend.Close()

	dispatcher, err := NewDispatcher(
		map[string]*Service{
			"default": &Service{
				Name:   "default",
				Origin: mustParseURL(defaultBackend.URL),
			},
			"foo": &Service{
				Name:   "foo",
				Origin: mustParseURL(fooBackend.URL),
			},
		},
		&Config{
			Rules: []ConfigRule{
				{
					ServiceName:     "default",
					HostPathMatcher: mustCompileHostPathMatcher("*/default/*"),
				},
				{
					ServiceName:     "foo",
					HostPathMatcher: mustCompileHostPathMatcher("*/foo/*"),
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	proxy := httptest.NewServer(NewProxyHandler(dispatcher))
	defer proxy.Close()

	reqGet := func(service, path string, status int) (*http.Response, error) {
		u := fmt.Sprintf("%s/%s%s", proxy.URL, service, path)
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", "testing")
		req.Header.Set("Status", strconv.Itoa(status))
		return http.DefaultClient.Do(req)
	}
	reqPost := func(service, path, body string, status int) (*http.Response, error) {
		u := fmt.Sprintf("%s/%s%s", proxy.URL, service, path)
		req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(body))
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", "testing")
		req.Header.Set("Status", strconv.Itoa(status))
		return http.DefaultClient.Do(req)
	}

	for _, service := range []string{"default", "foo"} {
		t.Run(service, func(t *testing.T) {
			t.Run("GET", func(t *testing.T) {
				res, err := reqGet(service, "/foo", 200)
				if err != nil {
					t.Error(err)
				}
				defer res.Body.Close()

				if res.StatusCode != 200 {
					t.Errorf("proxy status code should be 200 but got %d", res.StatusCode)
				}
				if s := res.Header.Get("Service"); s != service {
					t.Errorf("should proxy to %s, but got %s", service, s)
				}

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Error(err)
				}

				got := string(body)
				expected := fmt.Sprintf(heredoc.Doc(`
					GET /%s/foo
					Accept-Encoding: gzip
					Status: 200
					User-Agent: testing
					X-Forwarded-For: 127.0.0.1
				`), service)
				if diff := cmp.Diff(expected, got); diff != "" {
					t.Errorf("Unexpected response body: %s", got)
					t.Log(diff)
				}
			})

			t.Run("POST", func(t *testing.T) {
				res, err := reqPost(service, "/", "this is a body\n", 201)
				if err != nil {
					t.Error(err)
				}
				defer res.Body.Close()

				if res.StatusCode != 201 {
					t.Errorf("proxy status code should be 201 but got %d", res.StatusCode)
				}
				if s := res.Header.Get("Service"); s != service {
					t.Errorf("should proxy to default, but got %s", s)
				}

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Error(err)
				}

				got := string(body)
				expected := fmt.Sprintf(heredoc.Doc(`
					POST /%s/
					Accept-Encoding: gzip
					Content-Length: 15
					Status: 201
					User-Agent: testing
					X-Forwarded-For: 127.0.0.1
					this is a body
				`), service)
				if diff := cmp.Diff(expected, got); diff != "" {
					t.Errorf("Unexpected response body: %s", got)
					t.Log(diff)
				}
			})
		})
	}

	t.Run("X-Forwarded-For", func(t *testing.T) {
		u := fmt.Sprintf("%s%s", proxy.URL, "/default/bar")
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			t.Error(err)
		}

		req.Header.Set("User-Agent", "testing")
		req.Header.Set("Status", "200")
		req.Header.Set("X-Forwarded-For", "203.0.113.1")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			t.Errorf("proxy status code should be 200 but got %d", res.StatusCode)
		}
		if s := res.Header.Get("Service"); s != "default" {
			t.Errorf("should proxy to default, but got %s", s)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		got := string(body)
		expected := heredoc.Doc(`
			GET /default/bar
			Accept-Encoding: gzip
			Status: 200
			User-Agent: testing
			X-Forwarded-For: 127.0.0.1, 203.0.113.1
		`)
		if got != expected {
			t.Errorf("Unexpected response body: %s", got)
		}
	})

	t.Run("NoBackend", func(t *testing.T) {
		res, err := http.Get(proxy.URL)
		if err != nil {
			t.Error(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 404 {
			t.Errorf("proxy status code should be 404 but got %d", res.StatusCode)
		}
	})
}

func getBackendHandler(service string) http.Handler {
	replacer := strings.NewReplacer("\r\n", "\n")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Service", service)

		status := r.Header.Get("Status")
		if status != "" {
			s, err := strconv.Atoi(status)
			if err != nil {
				panic(err)
			}

			w.WriteHeader(s)
		}

		var buf bytes.Buffer
		r.Header.Write(&buf)
		headers := replacer.Replace(buf.String())

		io.WriteString(w, r.Method+" "+r.URL.RequestURI()+"\n")
		io.WriteString(w, headers)
		io.Copy(w, r.Body)
	})
}

type brokenResponseWriter struct{}

func (b brokenResponseWriter) Header() http.Header {
	return http.Header{}
}

func (b brokenResponseWriter) WriteHeader(status int) {
}

func (b brokenResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("broken")
}

func TestServiceProxyHandler(t *testing.T) {
	t.Run("FailedToCreateRequest", func(t *testing.T) {
		var reported []error
		reporter := ErrorReporterFunc(func(err error) {
			reported = append(reported, err)
		})

		handler := &serviceProxyHandler{service: &Service{Name: "default", Origin: mustParseURL("http://203.0.113.1")}, errorReporter: reporter}
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, &http.Request{
			Method: " INVALID ",
			URL:    mustParseURL("/"),
		})

		if len(reported) != 1 {
			t.Errorf("Unexpected reported errors: %v", reported)
		}

		result := recorder.Result()
		if result.StatusCode != http.StatusBadRequest {
			t.Errorf("Unexpected response status: %d", result.StatusCode)
		}
	})

	t.Run("FailedToRequestForBackend", func(t *testing.T) {
		var reported []error
		reporter := ErrorReporterFunc(func(err error) {
			reported = append(reported, err)
		})

		handler := &serviceProxyHandler{service: &Service{Name: "default", Origin: mustParseURL("http://203.0.113.1:99999")}, errorReporter: reporter}
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, &http.Request{
			Method: "GET",
			URL:    mustParseURL("/"),
		})

		if len(reported) != 1 {
			t.Errorf("Unexpected reported errors: %v", reported)
		}

		result := recorder.Result()
		if result.StatusCode != http.StatusBadGateway {
			t.Errorf("Unexpected response status: %d", result.StatusCode)
		}
	})

	t.Run("FailedToWriteResponse", func(t *testing.T) {
		defaultBackend := httptest.NewServer(getBackendHandler("default"))
		defer defaultBackend.Close()

		var reported []error
		reporter := ErrorReporterFunc(func(err error) {
			reported = append(reported, err)
		})

		handler := &serviceProxyHandler{service: &Service{Name: "default", Origin: mustParseURL(defaultBackend.URL)}, errorReporter: reporter}
		handler.ServeHTTP(brokenResponseWriter{}, &http.Request{
			Method: "GET",
			URL:    mustParseURL("/"),
		})

		if len(reported) != 1 {
			t.Errorf("Unexpected reported errors: %v", reported)
		} else if reported[0].Error() != "broken" {
			t.Errorf("Unexpected reported errors: %v", reported)
		}
	})
}

func TestErrorReporter(t *testing.T) {
	t.Run("Nop", func(t *testing.T) {
		nopErrorReporter.ReportError(errors.New("foo"))
	})
}

func TestCreateProxyRequest(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		handler := &serviceProxyHandler{service: &Service{Name: "default", Origin: mustParseURL("http://localhost:8080")}}
		_, err := handler.createProxyRequest(&http.Request{
			Method: " INVALID ",
			URL:    mustParseURL("http://localhost:8080/"),
		})
		if err == nil {
			t.Errorf("should be error")
		}
	})
}

func TestFilterHeaders(t *testing.T) {
	t.Run("Keep", func(t *testing.T) {
		h := http.Header{}
		h.Set("Server", "foo")
		filterHeaders(h)

		expected := http.Header{}
		expected.Set("Server", "foo")
		if diff := cmp.Diff(expected, h); diff != "" {
			t.Errorf("should not be filterd any fields but got %s", diff)
		}
	})

	t.Run("Connection", func(t *testing.T) {
		h := http.Header{}
		h.Set("Server", "foo")
		h.Set("Connection", "Keep-Alive")
		h.Set("Keep-Alive", "timeout=5, max=1000")
		filterHeaders(h)

		expected := http.Header{}
		expected.Set("Server", "foo")
		if diff := cmp.Diff(expected, h); diff != "" {
			t.Errorf("should be filterd some fields but got %s", diff)
		}
	})
}

func TestGetRemoteIP(t *testing.T) {
	if ip := getRemoteIP(&http.Request{RemoteAddr: "203.0.113.1"}); ip != "203.0.113.1" {
		t.Errorf("should be 203.0.113.1 but got %s", ip)
	}

	if ip := getRemoteIP(&http.Request{RemoteAddr: "203.0.113.1:12345"}); ip != "203.0.113.1" {
		t.Errorf("should be 203.0.113.1 but got %s", ip)
	}
}
