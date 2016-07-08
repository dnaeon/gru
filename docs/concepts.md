## Concepts

Gru is designed around the following concepts, each of which is
explained below.

## Resource

Resources are the core components in Gru. Each resource is
responsible for handling a particular task in an idempotent manner, e.g.
management of packages, management of services, executing commands, etc.

## Module

A module is essentially a [Lua](https://www.lua.org/) module.

Lua is used to form the foundation of the DSL language used in Gru.

Within modules resources are being created and registered to the
catalog.

## Catalog

The catalog represents a collection of resources, which were
created after evaluating a given module.

Before processing the catalog all resources are first
[topologically sorted](https://en.wikipedia.org/wiki/Topological_sorting),
in order to determine the proper order of evaluation and processing.

## Task

A task represents a message to remote minions, that a given
module should be evaluated and result should be returned.

The task itself also bundles additional meta data, such as the
unique id of the task, the time when task has been received,
processed, etc.

Tasks are sent out to minions using [etcd](https://github.com/coreos/etcd).

## Client

The client is used to push tasks to minions, retrieve results,
generate reports about minions, etc. It is the frontend application of
Gru.

## Minion

The minion is a remote system which receives and processes tasks from
clients.
