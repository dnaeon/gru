## Gru - A Simple Orchestration Framework

[![Build Status](https://travis-ci.org/dnaeon/gru.svg)](https://travis-ci.org/dnaeon/gru)
[![GoDoc](https://godoc.org/github.com/dnaeon/gru?status.svg)](https://godoc.org/github.com/dnaeon/gru)
[![Go Report Card](https://goreportcard.com/badge/github.com/dnaeon/gru)](https://goreportcard.com/report/github.com/dnaeon/gru)

Gru is a simple orchestration framework written in Go, which
allows you to manage your UNIX/Linux systems with ease.

## Installation

The easiest way to get Gru is to use one of the pre-built binaries
from the [Gru releases page](https://github.com/dnaeon/gru/releases/).

For those wanting to try out the latest version of Gru you should
follow these instructions instead:

```bash
$ git clone https://github.com/dnaeon/gru
$ cd gru
$ make
```

Gru uses
[etcd](https://github.com/coreos/etcd) for discovery of minions and
communication between the minions and clients. In order to use Gru,
first make sure that you have `etcd` up and running.

## Usage

[![asciicast](https://asciinema.org/a/35920.png)](https://asciinema.org/a/35920)

## Status

Experimental. Gru is in early development stage. Consider the
API unstable for as now as things change rapidly.

## License

Gru is Open Source and licensed under the
[BSD License](http://opensource.org/licenses/BSD-2-Clause)
