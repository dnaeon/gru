## Installation

Gru is being built and tested against Golang tip. You will need
Golang tip in order to build and run Gru.

The easiest way to get the tip version of Golang is by using
[gimme](https://github.com/travis-ci/gimme), e.g.

```bash
$ gimme tip
$ source ~/.gimme/envs/gotip.env
```

Gru is also using the
[libgit2/git2go bindings](https://github.com/libgit2/git2go), which
require version `0.24.0` of `libgit2` to be installed. You could
install `libgit2` using your package manager or use the
[contrib/misc/build-libgit2.sh](../contrib/misc/build-libgit2.sh)
script if your package manager doesn't provide version 0.24.0 of libgit2.

[etcd](https://github.com/coreos/etcd) is used for discovery of
minions and communication between the minions and clients, so before
you can use Gru you need to make sure that you have `etcd` up and
running.

For more information about installing and configuring
[etcd](https://github.com/coreos/etcd), please refer to the
[official etcd documentation](https://coreos.com/etcd/docs/latest/).

Once you've got all requirements installed, installing Gru is as
easy as executing these commands below.

```bash
$ git clone https://github.com/dnaeon/gru
$ cd gru
$ make
```

## Tests

You can run the tests by executing the command below.

```bash
$ make test
```
