package gaedispemu

import (
	"io"
	"net/http"
	"strings"
)

// ErrorReporter is error reporter interface for proxy handler
type ErrorReporter interface {
	ReportError(error)
}

type nopErrorReporterType struct{}

var nopErrorReporter = nopErrorReporterType{}

func (r nopErrorReporterType) ReportError(err error) {}

// ErrorReporterFunc is function interface for ErrorReporter
type ErrorReporterFunc func(error)

// ReportError calls self as a function
func (r ErrorReporterFunc) ReportError(err error) {
	r(err)
}

// NewProxyHandler creates a new proxy handler
func NewProxyHandler(dispatcher Dispatcher) http.Handler {
	return &proxyHandler{dispatcher: dispatcher, errorReporter: nopErrorReporter}
}

// NewProxyHandlerWithReporter creates a new proxy handler with error reporter
func NewProxyHandlerWithReporter(dispatcher Dispatcher, errorReporter ErrorReporter) http.Handler {
	return &proxyHandler{dispatcher: dispatcher, errorReporter: errorReporter}
}

type proxyHandler struct {
	dispatcher    Dispatcher
	errorReporter ErrorReporter
}

var _ http.Handler = (*proxyHandler)(nil)

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service := h.dispatcher.Dispatch(r.URL.Host, r.URL.Path)
	if service == nil {
		http.Error(w, "No such backend for the URL: "+r.URL.Path, http.StatusNotFound)
		return
	}

	next := &serviceProxyHandler{service: service, errorReporter: h.errorReporter}
	next.ServeHTTP(w, r)
}

// SEE ALSO: RFC2616
var nopHeadersByHop = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

type serviceProxyHandler struct {
	service       *Service
	errorReporter ErrorReporter
}

var _ http.Handler = (*serviceProxyHandler)(nil)

func (h *serviceProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := h.createProxyRequest(r)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusBadRequest)
		h.errorReporter.ReportError(err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to request for backend", http.StatusBadGateway)
		h.errorReporter.ReportError(err)
		return
	}

	err = h.proxyResponse(w, res)
	if err != nil {
		h.errorReporter.ReportError(err)
	}
}

func (h *serviceProxyHandler) createProxyRequest(src *http.Request) (*http.Request, error) {
	u := h.service.Origin.ResolveReference(src.URL)
	dst, err := http.NewRequest(src.Method, u.String(), src.Body)
	if err != nil {
		return nil, err
	}

	// proxy headers
	copyHeader(dst.Header, src.Header)
	dst.Header.Set("X-Forwarded-For", getNewForwardedIPs(src))
	filterHeaders(dst.Header)
	dst.ContentLength = src.ContentLength

	return dst, nil
}

func (h *serviceProxyHandler) proxyResponse(w http.ResponseWriter, res *http.Response) error {
	defer res.Body.Close()

	// proxy headers
	filterHeaders(res.Header)
	copyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)

	// proxy body
	_, err := io.Copy(w, res.Body)
	return err
}

func filterHeaders(h http.Header) {
	tokens := h[http.CanonicalHeaderKey("Connection")]
	if tokens != nil && len(tokens) > 0 {
		for _, token := range tokens {
			parts := strings.Split(token, ",")
			for _, part := range parts {
				trimed := strings.TrimSpace(part)
				h.Del(trimed)
			}
		}
	}

	for _, key := range nopHeadersByHop {
		h.Del(key)
	}
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		dst[key] = values
	}
}

func getNewForwardedIPs(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded == "" {
		return getRemoteIP(r)
	}

	remoteIP := getRemoteIP(r)
	return remoteIP + ", " + forwarded
}

func getRemoteIP(r *http.Request) string {
	// remove port
	index := strings.Index(r.RemoteAddr, ":")
	if index == -1 {
		return r.RemoteAddr
	}

	return r.RemoteAddr[:index]
}
