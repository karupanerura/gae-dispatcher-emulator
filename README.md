# gae-dispatcher-emulator

[![Build Status](https://travis-ci.org/karupanerura/gae-dispatcher-emulator.svg?branch=master)](https://travis-ci.org/karupanerura/gae-dispatcher-emulator)
[![codecov](https://codecov.io/gh/karupanerura/gae-dispatcher-emulator/branch/master/graph/badge.svg)](https://codecov.io/gh/karupanerura/gae-dispatcher-emulator)
[![GoDoc](https://godoc.org/github.com/karupanerura/gae-dispatcher-emulator?status.svg)](http://godoc.org/github.com/karupanerura/gae-dispatcher-emulator)

Unofficial Google App Engine Dispatcher Emulator

## Description

`gae-dispatcher-emulator` is an unofficial emulator for `Google App Engine` dispatcher service.
This works like a local reverse proxy, and it behave by `dispatch.yaml` or `dispatch.xml`.

```console
$ (cd default; dev_appserver.py --port=8081 | tee -a dev.log) &
$ (cd mobile-backend; dev_appserver.py --port=8082 | tee -a dev.log) &
$ (cd static-backend; dev_appserver.py --port=8083 | tee -a dev.log) &
$ gae-dispatcher-emulator -c dispatch.yaml -s default:localhost:8081 -s mobile-frontend:localhost:8082 -s static-backend:localhost:8083
```

I suggest to use it with [foreman](http://ddollar.github.io/foreman/) to launch/shutdown services consistently.

## Installation

```bash
go get -u github.com/karupanerura/gae-dispatcher-emulator/...
```

## Usage

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
