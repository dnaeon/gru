## Concepts

Gru is designed around the following concepts, each of which is
explained below.

## Resource

Resources are the core elements in Gru modules. Each resource is
responsible for handling a particular task, e.g. manage package
installations, enabling and starting of services, etc.

Resources are evaluated and processed by minions and are
idempotent.

## Module

A module is a collection of resources.

A module can import other modules, thus allowing for better
organization of logic and code re-use.

Modules are what forms the DSL language of Gru and can be written in
either [HCL](https://github.com/hashicorp/hcl) or
[JSON](http://www.json.org/).

## Catalog

The catalog is a collection of modules. Prior to catalog processing
the resources from all modules in the catalog are
[topologically sorted](https://en.wikipedia.org/wiki/Topological_sorting),
in order to determine the proper order of evaluation and processing of
the resources.

## Task

A task bundles a catalog with some meta data such as the
unique id of the task, the time when task has been received,
processed, etc.

Tasks are received and processed by minions.

## Client

The Gru client is used to push tasks to minions, retrieve results,
generate reports about minions, etc. It is the frontend application of
Gru.

## Minion

The minion is a remote system which receives and processes
tasks from clients.
