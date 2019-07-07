# gae-dispatcher-emulator

[![Build Status](https://travis-ci.org/karupanerura/gae-dispatcher-emulator.svg?branch=master)](https://travis-ci.org/karupanerura/gae-dispatcher-emulator)
[![codecov](https://codecov.io/gh/karupanerura/gae-dispatcher-emulator/branch/master/graph/badge.svg)](https://codecov.io/gh/karupanerura/gae-dispatcher-emulator)
[![GoDoc](https://godoc.org/github.com/karupanerura/gae-dispatcher-emulator?status.svg)](http://godoc.org/github.com/karupanerura/gae-dispatcher-emulator)

Unofficial Google App Engine Dispatcher Emulator

## INSTALL

```
go get github.com/karupanerura/gae-dispatcher-emulator/...
```

```
Usage:
  gae-dispatcher-emulator [OPTIONS]

Application Options:
  -c, --config=	 dispatch.xml or dispatch.yaml
  -s, --service= service map (e.g. --service default:localhost:8081 --service admin:localhost:8082)
  -l, --listen=	 listening host:port (localhost:3000 is default) (default: localhost:3000)
  -v, --verbose	 verbose output for proxy request

Help Options:
  -h, --help	 Show this help message
```
