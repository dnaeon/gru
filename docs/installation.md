## Installation

In order to build Gru you need Go version 1.7 or later.

Building Gru is as easy as executing the commands below.

```bash
$ git clone https://github.com/dnaeon/gru
$ cd gru
$ make
```

## Optional requirements

The optional requirements listed below are needed if you need to
orchestrate remote systems. They are not required if you use Gru in
stand-alone mode.

[etcd](https://github.com/coreos/etcd) is used for discovery of
minions and communication between the minions and clients, so
in order to orchestrate remote minions you need to make sure that you
have `etcd` up and running, so that remote minions can connect to it.

For more information about installing and configuring
[etcd](https://github.com/coreos/etcd), please refer to the
[official etcd documentation](https://coreos.com/etcd/docs/latest/).

[Git](https://git-scm.com/) is used for syncing code and data
files to the remote minions, so make sure that you have Git
installed as well.

## Shell Completion

You can enable shell autocompletion by sourcing the
correct file from the [contrib/autocomplete](../contrib/autocomplete)
directory for your shell.

For instance to enable bash autocompletion on an Arch Linux system,
you would do.

```bash
$ sudo cp contrib/autocomplete/bash_autocomplete /usr/share/bash-completion/completions/gructl
```

Note that you will need to install the `bash-completion`
package, if you don't have it installed already.

## Tests

You can run the tests by executing the command below.

```bash
$ make test
```
