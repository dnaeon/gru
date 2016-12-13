## Gru - Orchestration made easy with Go and Lua

[![Build Status](https://travis-ci.org/dnaeon/gru.svg)](https://travis-ci.org/dnaeon/gru)
[![GoDoc](https://godoc.org/github.com/dnaeon/gru?status.svg)](https://godoc.org/github.com/dnaeon/gru)
[![Go Report Card](https://goreportcard.com/badge/github.com/dnaeon/gru)](https://goreportcard.com/report/github.com/dnaeon/gru)
[![Join the chat at https://gitter.im/dnaeon/gru](https://badges.gitter.im/dnaeon/gru.svg)](https://gitter.im/dnaeon/gru?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Codewake](https://www.codewake.com/badges/ask_question.svg)](https://www.codewake.com/p/gru)

Gru is a fast and concurrent orchestration framework powered
by Go and Lua, which allows you to manage your UNIX/Linux systems
with ease.

## Documentation

You can find the latest documentation [here](docs/).

The API documentation is available [here](https://godoc.org/github.com/dnaeon/gru).

## Features

* Written in fast, compiled language - [Go](https://golang.org/)
* Uses a fast, lightweight, embeddable, scripting
  language as the DSL - [Lua](https://www.lua.org/)
* Concurrent execution of idempotent operations
* Distributed - using [etcd](https://github.com/coreos/etcd) for node
  discovery and communication and
  [Git](https://git-scm.com/) for version control and data sync
* Easy to deploy - comes with a single, statically linked binary
* Suitable for orchestration and configuration management

## Status

Gru is in constant development. Consider the API unstable as
things may change without a notice.

## Contributions

Gru is hosted on [Github](https://github.com/dnaeon/gru).
Please contribute by reporting issues, suggesting features or by
sending patches using pull requests.

## License

Gru is Open Source and licensed under the
[BSD License](http://opensource.org/licenses/BSD-2-Clause).

## References

References to articles related to this project in one way or another.

* [Managing VMware vSphere environment with Go and Lua by using Gru orchestration framework](http://dnaeon.github.io/gru-vmware-vsphere-mgmt/)
* [Introducing triggers in Gru orchestration framework](http://dnaeon.github.io/introducing-triggers-in-gru/)
* [Puppet vs Gru - Benchmarking Speed & Concurrency](http://dnaeon.github.io/puppet-vs-gru-benchmarking-speed-and-concurrency/)
* [Extending Lua with Go types](http://dnaeon.github.io/extending-lua-with-go-types/)
* [Choosing Lua as the data description and configuration language](http://dnaeon.github.io/choosing-lua-as-the-ddl-and-config-language/)
* [Creating an orchestration framework in Go](http://dnaeon.github.io/gru-orchestration-framework/)
* [Dependency graph resolution algorithm in Go](http://dnaeon.github.io/dependency-graph-resolution-algorithm-in-go/)
* [Orchestration made easy with Gru v0.2.0](http://dnaeon.github.io/orchestration-made-easy-with-gru-v0.2.0/)
* [Membership test in Go](http://dnaeon.github.io/membership-test-in-go/)
* [Testing HTTP interactions in Go](http://dnaeon.github.io/testing-http-interactions-in-go/)
* [Concurrent map and slice types in Go](http://dnaeon.github.io/concurrent-maps-and-slices-in-go/)
* [Lua as a Configuration And Data Exchange Language](https://www.netbsd.org/~mbalmer/lua/lua_config.pdf)
