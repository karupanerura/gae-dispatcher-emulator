package gaedispemu

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
)

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
				if got != expected {
					t.Errorf("Unexpected response body: %s", got)
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
					Status: 201
					User-Agent: testing
					X-Forwarded-For: 127.0.0.1
					this is a body
				`), service)
				if got != expected {
					t.Errorf("Unexpected response body: %s", got)
				}
			})
		})
	}

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
