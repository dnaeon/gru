## Installation

Gru is being built and tested against the Go tip version.

In order to build Gru, simply follow the steps below.

```bash
$ git clone https://github.com/dnaeon/gru
$ cd gru
$ make
```

Gru uses
[etcd](https://github.com/coreos/etcd) for discovery of minions and
communication between the minions and clients, so before you can use
Gru you need to make sure that you have `etcd` up and running.

For installing and configuring [etcd](https://github.com/coreos/etcd),
please refer to the
[official etcd documentation](https://coreos.com/etcd/docs/latest/).

## Tests

You can run the tests by executing the command below:

```bash
$ make test
```
