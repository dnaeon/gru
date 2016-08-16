## Installation

In order to build Gru you will need Go version 1.7 or later.

[etcd](https://github.com/coreos/etcd) is used for discovery of
minions and communication between the minions and clients, so before
you can use Gru you need to make sure that you have `etcd` up and
running.

For more information about installing and configuring
[etcd](https://github.com/coreos/etcd), please refer to the
[official etcd documentation](https://coreos.com/etcd/docs/latest/).

[Git](https://git-scm.com/) is used for syncing code and data
files to the remote minions, so make sure that you have Git
installed as well.

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
