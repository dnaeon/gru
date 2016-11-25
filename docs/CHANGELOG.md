## 0.5.0 (November 24, 2016)

* Added support for managing VMware vSphere environments. Check the
  list of new resources below.
* Added `Property` interface type used for evaluating and setting
  properties of resources.
* Impelemented resource triggers. Example code using triggers can be
  found [here](https://github.com/dnaeon/gru/blob/master/site/code/triggers.lua)
* Existing resources have been refactored to use properties
* Use latest version of `go-vcr` for the integration tests.
* Use ldflags when building the `gructl` binary to reduce the resulting binary size
* Start using `context.Context` from Go 1.7 instead of `golang.org/x/net/context`
* The `resource.Logf()` function is now registered in Lua and
  available from the `stdlib` table
* Go functions can be registered in Lua as well allowing for further
  extending the Lua-based DSL with Go functionality

Changes in `resource.Resource` interface type:

* The `Update()` method is now gone. From now on we use properties to
  describe changes that can be evaluated and set on resources.
* The `Log()` method is gone, instead use `resource.Logf()`
* New method `Properties()` has been added, which returns the list of
  properties for a resource.
* New method `Initialize()` has been added, which is used for
  performing initialization of the resource if needed, e.g.
  establishing connection to a remote API endpoint.
* New method `Close()` has been added, which is called after a resource
  has been processed can be used for cleanup tasks.
* Resources for which properties make no sense when absent now
  return the `resource.ErrResourceAbsent` error.
* New method `SubscribedTo` returns a `resource.TriggerMap`
  containing the resources to which it is subscribed and the Lua
  function that will be executed if the resource has changed.

New resources implemented:

* `File` resource now manages only files instead of files and directories.
* `Directory` manages directories.
* `Link` manages symbolic and hard file links.
* `SysRc` resource manages sysrc variables on FreeBSD systems. Contributed by @kaey
* `PkgNg` resource manages packages on FreeBSD systems. Contributed by @kaey
* `Service` resource for FreeBSD systems. Contributed by @kaey
* `Datacenter` resource manages datacenters in a VMware vSphere environment
* `Cluster` resource manages VMware vSphere clusters
* `ClusterHost` manages ESXi hosts in a VMware vSphere cluster.
* `Host` manages ESXi host settings in a VMware vSphere environment
* `VirtualMachine` resource manages virtual machines in a VMware vSphere environment
* `DatastoreNfs` resource manages NFS datastores on ESXi hosts

Changes in the `catalog` package.

* Added `catalog.Status` type
* `Catalog.Run()` now returns `catalog.Status` instance, so that
  result can be further examined if needed
* Resource properties are processed after executing the resource

## 0.4.0 (August 24, 2016)

* Use Go 1.7 as the default stable version used for building and testing
* Added support for concurrent resource processing (!)
* Suitable resources are now scheduled for concurrent execution
* Created [Gitter channel for Gru](https://gitter.im/dnaeon/gru)
* Created [Codewake channel for Gru](https://www.codewake.com/p/gru)
* Resources can define their own custom `present` and `absent` states
* Added new required methods for `Resource` interface -
  `Validate()`, `IsConcurrent()`, `GetPresentStates()` and `GetAbsentStates()`
* Renamed `Resource.BaseResource` type as `Resource.Base`
* Initial implementation of `utils.List` and `utils.String` types used to
  provide membership test operations
* The `Update` field of `resource.State` type been renamed to `Outdated`
* Resource processing logic has been simplified in the `catalog` package
* Added tests for the `catalog` package
* Support only direct dependencies by using `require` - `before` and `after` are gone
* Display status of applied resource after a catalog run
* During a catalog run resources which have failed dependencies are now skipped
* `Catalog.Add()` registeres only non-nil resources
* Added `Mute` field to `resource.Shell` type, which suppresses output from shell commands
* Be able to set concurrency level when applying configuration with `gructl apply` and `gructl serve`
* Implemented `Reversed()` method on `graph.Graph` type
* Added support for shell autocompletion
* Added support for resource namespaces in Lua
* `gructl graph` now generates the reversed graph of resources as well
* Updated documentation

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
