## 0.2.0 (May 27, 2016)

* Introduced [resources](resources/) as a way to define and execute idempotent operations
* Integrated [HCL](https://github.com/hashicorp/hcl) as a way to express configuration
* Initial implementation of the [module](module/) package
* Initial implementation of the [graph](graph/) package and resource dependencies
* Initial implementation of the [catalog](catalog/) package
* Initial implementation of the `site repo`, which allows
  remote minions to sync data and modules from an upstream Git repository
* Minions now understand the notion of task environments, which
  essentially are Git branches in the `site repo`
* Initial implementation of some core resources - `package`, `service`,
  `file`, `shell`, `pacman` and `yum`
* Improved the existing documentation and created the
  [quickstart guide](docs/quickstart.md) guide
* supervisord and systemd service configurations created and available
  from the [contrib](contrib/) directory
* And many other small fixes and improvements

## 0.1.0 (Feb 17, 2016)

Initial release of Gru.

* Initial implementation of the core interfaces based on [etcd](https://github.com/coreos/etcd)

[![asciicast](https://asciinema.org/a/35920.png)](https://asciinema.org/a/35920)
