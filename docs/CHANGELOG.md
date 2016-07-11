## 0.3.0 (July 11, 2016)

* [Lua](https://www.lua.org/) has been integrated as the DSL language used by Gru.
* [HCL](https://github.com/hashicorp/hcl) has been removed as we now use Lua.
* `resource.LuaRegisterBuiltin` was implemented for registering resources into Lua.
* `resource.Provider` type signature has changed after the adoption of Lua.
* Renamed `ResourceID()` method for `resource.Resource` to `ID()`.
* `SetType()` is no longer a method of the `resource.Resource` interface.
* Removed the `Title` field from `resource.BaseResource` type.
* Services managed by `resource.Service` are now enabled by default.
* The `module.Collection` type has moved to `resource.Collection`.
* Deprecated the `module` package, which used to implement modules using HCL.
  Since we have moved to Lua, there is no need of this package anymore.
* Initial implementation of `resource.DefaultConfig` which provides configuration
  settings for all resources. `resource.DefaultConfig` is being injected by
  external packages with the proper configuration.
* Deprecated `gructl module` command.
* Deprecated `gructl validate` command.
* Deprecated `gructl resource` command.
* Switched to `urfave/cli` instead of `codegangsta/cli`
* The `log` function is now registered to Lua for logging events.
* Dropped requirement for `libgit2`. Minions now sync the
  `site repo` using the `utils.GitRepo` implementation instead.
* Added tests for all resources
* The resource documentation and examples are now available in
  [godoc](https://godoc.org/github.com/dnaeon/gru).
* Updated documentation and the quickstart guide to reflect the
  adoption of Lua as the DSL language.

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
