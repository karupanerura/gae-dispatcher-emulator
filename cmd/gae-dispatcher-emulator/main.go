// Usage:
//   gae-dispatcher-emulator [OPTIONS]
//
// Application Options:
//   -c, --config=  dispatch.xml or dispatch.yaml
//   -s, --service= service map (e.g. --service default:localhost:8081 --service admin:localhost:8082)
//   -l, --listen=  listening host:port (localhost:3000 is default) (default: localhost:3000)
//   -v, --verbose  verbose output for proxy request
//
// Help Options:
//   -h, --help     Show this help message
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	gaedispemu "github.com/karupanerura/gae-dispatcher-emulator"
	"github.com/motemen/go-loghttp"
)

type options struct {
	ConfigFile string   `short:"c" long:"config" description:"dispatch.xml or dispatch.yaml" required:"true"`
	Services   []string `short:"s" long:"service" description:"service map (e.g. --service default:localhost:8081 --service admin:localhost:8082)" required:"true"`
	ListenAddr string   `short:"l" long:"listen" description:"listening host:port (localhost:3000 is default)" default:"localhost:3000"`
	Verbose    bool     `short:"v" long:"verbose" description:"verbose output for proxy request"`
}

func main() {
	var opts options
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Printf("Failed to parse args: %v", err)
		os.Exit(1)
	}

	handler, err := createProxyHandler(&opts)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}

	if opts.Verbose {
		http.DefaultTransport = loghttp.DefaultTransport
	}

	server := opts.getServer(handler)
	log.Printf("Listen on %s", opts.ListenAddr)
	log.Fatal(server.ListenAndServe())
}

type loggingErrorReporter struct{}

func (r loggingErrorReporter) ReportError(err error) {
	log.Printf("ERROR: %v", err)
}

func createProxyHandler(opts *options) (http.Handler, error) {
	loader := opts.getConfigLoader()
	if loader == nil {
		return nil, fmt.Errorf("Failed to determine config type for %q", opts.ConfigFile)
	}

	config, err := loader.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %v", err)
	}

	services, err := opts.getServicsMap()
	if err != nil {
		return nil, err
	}

	dispatcher, err := gaedispemu.NewDispatcher(services, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to mapping backend: %v", err)
	}

	reporter := loggingErrorReporter{}
	return gaedispemu.NewProxyHandlerWithReporter(dispatcher, reporter), nil
}

func (o options) getConfigLoader() gaedispemu.ConfigLoader {
	if strings.HasSuffix(o.ConfigFile, ".xml") {
		return gaedispemu.NewXMLConfigLoader(o.ConfigFile)
	} else if strings.HasSuffix(o.ConfigFile, ".yaml") {
		return gaedispemu.NewXMLConfigLoader(o.ConfigFile)
	}

	return nil
}

func (o options) getServicsMap() (map[string]*gaedispemu.Service, error) {
	m := make(map[string]*gaedispemu.Service, len(o.Services))
	for _, service := range o.Services {
		index := strings.Index(service, ":")
		if index == -1 {
			return nil, fmt.Errorf("Invalid service map format: %s", service)
		}

		name := service[:index]
		if _, ok := m[name]; ok {
			return nil, fmt.Errorf("Duplicated service name: %s", name)
		}

		origin, err := parseOrigin(service[index+1:])
		if err != nil {
			return nil, fmt.Errorf("Invalid service map format: %s (%v)", service, err)
		}
		m[name] = &gaedispemu.Service{
			Name:   name,
			Origin: origin,
		}
	}
	return m, nil
}

func parseOrigin(s string) (*url.URL, error) {
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return url.Parse(s)
	}

	return url.Parse("http://" + s)
}

func (o options) getServer(h http.Handler) *http.Server {
	return &http.Server{
		Addr:     o.ListenAddr,
		Handler:  h,
		ErrorLog: log.New(os.Stderr, "", log.LstdFlags),
	}
}
