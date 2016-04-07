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
communication between the minions and clients, so before you can use
Gru you need to make sure that you have `etcd` up and running.

For installing and configuring [etcd](https://github.com/coreos/etcd),
please refer to the
[official etcd documentation](https://coreos.com/etcd/docs/latest/).

## Tests

You can also run the tests by executing the command below:

```bash
$ make test
```
