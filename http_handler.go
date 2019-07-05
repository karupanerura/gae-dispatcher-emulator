package gaedispemu

import (
	"io"
	"net/http"
	"strings"
)

// NewProxyHandler creates a new proxy handler
func NewProxyHandler(dispatcher Dispatcher) http.Handler {
	return &dispatcherHandler{dispatcher: dispatcher}
}

type dispatcherHandler struct {
	dispatcher Dispatcher
}

var _ http.Handler = (*dispatcherHandler)(nil)

func (h *dispatcherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service := h.dispatcher.Dispatch(r.URL.Host, r.URL.Path)
	if service == nil {
		http.Error(w, "No such backend for the URL: "+r.URL.Path, http.StatusNotFound)
		return
	}

	next := &serviceProxyHandler{service: service}
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
	service *Service
}

var _ http.Handler = (*serviceProxyHandler)(nil)

func (h *serviceProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := h.createProxyRequest(r)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusBadRequest)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to request for proxy", http.StatusInternalServerError)
		return
	}

	h.proxyResponse(w, res)
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

	return dst, nil
}

func (h *serviceProxyHandler) proxyResponse(w http.ResponseWriter, res *http.Response) {
	defer res.Body.Close()

	// proxy headers
	filterHeaders(res.Header)
	copyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)

	// proxy body
	_, err := io.Copy(w, res.Body)
	if err != nil {
		// TODO
	}
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
	addr := r.RemoteAddr

	// remove port
	index := strings.Index(addr, ":")
	if index != -1 {
		addr = r.RemoteAddr[:index]
	}

	return addr
}
