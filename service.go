package gaedispemu

import "net/url"

// Service is a GAE service and backend origin
type Service struct {
	Name   string
	Origin *url.URL
}
