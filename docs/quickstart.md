## Quickstart

Welcome to Gru!

Considering that you have installed and configured Gru, we will now
walk you through the creation of a simple module, which will
take care of installing and configuring
[memcached](https://memcached.org/) for us.

The instructions here should be pretty simple for everyone to follow,
and should serve as an example on creating modules for Gru.

The instructions presented in this document have been tested on an
[Arch Linux](https://www.archlinux.org/) system, but they should also
apply to other systems for which Gru has support.

## Setting up the site repo

The `site repo` in Gru is what contains modules and data files.

Technically speaking the site repo is a Git repository,
which is being distributed to our minions, this way making it possible
for remote systems to sync from it, load modules and later on
process them.

This is how a typical site repo structure looks like.

```bash
$ tree site
site
├── data
└── modules

2 directories, 0 files
```

The `modules` directory is where Gru modules reside, while `data`
directory is being used for static content and file templates.

You can also find an example site repo in the
[example site repo](../site) directory from the Gru repository.

This is also the place where you can find the modules we will
prepare in this document, so you might want to grab
the example site repo first while you work on the instructions from
this document.

## Writing the module

Modules in Gru are expressed in
[HCL](https://github.com/hashicorp/hcl) and reside in the
`modules` directory of the site repo as already mentioned above.

In the beginning of this document we have mentioned that we will be
installing and configuring [memcached](https://memcached.org/) on our
systems. The steps we need to perform in order to do that can be
summarized as installing the needed package, configuring the service
and afterwards starting the service.

First, let's begin with installing the requred packages by creating
our first resource.

```hcl
package "memcached" {
  state = "present"
}
```

By default the `memcached` service listens only on localhost, so
if we want to change that and listen on all interfaces we will have to
adjust the
[systemd](https://www.freedesktop.org/wiki/Software/systemd/) unit
for the service.

One way to achieve that is to use systemd drop-in units, and that is
what we will do now.

First, let's create the needed directory for our drop-in unit.

```hcl
file "/etc/systemd/system/memcached.service.d" {
  state = "present"
  filetype = "directory"
  require = [
    "package[memcached]",
  ]
}
```

Now, let's install the actual drop-in unit file.

```hcl
file "/etc/systemd/system/memcached.service.d/override.conf" {
  state = "present"
  mode = 0644
  source = "data/memcached/memcached-override.conf"
  require = [
    "file[/etc/systemd/system/memcached.service.d]",
  ]
}
```

The above `file` resource will take care of installing the
`override.conf` drop-in unit to it's correct location.

You may have also noticed the `source` parameter that we use in our
resource - that parameter tells Gru where in the site directory is the
actual source file located.

And this is what the actual drop-in unit file looks like, which
will be used as the source for our resource.

```ini
[Service]
ExecStart=
ExecStart=/usr/bin/memcached
```

Once we install the systemd drop-in unit we need to tell
`systemd(1)` to re-read it's configuration, so the next resource
takes care of that as well.

```hcl
shell "systemctl daemon-reload" {
  require = [
    "file[/etc/systemd/system/memcached.service.d/override.conf]",
  ]
}
```

And finally let's create a resource that enables and starts the
memcached service.

```hcl
service "memcached" {
  state = "running"
  enable = true
  require = [
    "package[memcached]",
    "file[/etc/systemd/system/memcached.service.d/override.conf]",
  ]
}
```

With all that we have now created our first module, which should
take care of installing and configuring memcached for us!

One last step that we should also consider doing is to validate our
configuration. In order to do that we will use the `gructl validate`
command.

```bash
$ gructl validate memcached
Loaded 5 resources from 1 modules
```

We can also see our new module being successfully discovered
using the `gructl module` command.

```bash
$ gructl module memcached
MODULE          PATH
memcached       .../gru/site/modules/memcached.hcl
```

And this is how our site repo looks like once we have everything in
place.

```bash
$ tree site/
site/
├── data
│   └── memcached
│       └── memcached-override.conf
└── modules
    └── memcached.hcl

3 directories, 2 files
```

## Resource Dependencies

In the previous chapter of this document we have created a number of
resources, which took care of installing and configuring memcached.

What you should have also noticed is that in most of the resources we
have used these `require` parameters.

The `require` parameter is used for creating resource dependencies.

Before the resources are being processed by the `catalog`, Gru is
building a [DAG graph](https://en.wikipedia.org/wiki/Directed_acyclic_graph)
of all resources and attempts to perform a
[topological sort](https://en.wikipedia.org/wiki/Topological_sorting)
on them, in order to determine the proper ordering of resource execution.

Considering the example memcached module we have created in the
previous chapter, let's see what it's DAG graph looks like.

In order to do that we will use the `gructl graph` command.

```dot
digraph resources {
        nodesep=1.0
        node [shape=box]
        edge [style=dashed]
        "service[memcached]" -> "package[memcached]"
        "service[memcached]" -> "file[/etc/systemd/system/memcached.service.d/override.conf]"
        "shell[systemctl daemon-reload]" -> "file[/etc/systemd/system/memcached.service.d/override.conf]"
        "file[/etc/systemd/system/memcached.service.d]" -> "package[memcached]"
        "file[/etc/systemd/system/memcached.service.d/override.conf]" -> "file[/etc/systemd/system/memcached.service.d]"
}
```

The result of the `gructl graph` command is a representation of the
dependency graph in the
[DOT language](https://en.wikipedia.org/wiki/DOT_(graph_description_language)).

If we pipe the above result to `dot(1)` we can generate a visual
representation of our graph, e.g.:

```bash
$ gructl graph --resources memcached | dot -O -Tpng
```

And this is how the dependency graph for our memcached module looks like.

![memcached graph dependencies](images/memcached-dag.png)

Using `gructl graph` we can see what the resource execution
ordering would look like and it can also help us identify
circular dependencies in our resources and modules.

## Applying Configuration

Applying configuration in Gru can be done in a couple of ways -
standalone and orchestration mode.

In standalone mode configuration is being applied on the local system,
while in orchestration mode a task is being pushed to the remote
minions for processing.

### Standalone mode

Let's see how we can apply the configuration from the module we've
prepared so far on the local system using the standalone mode.

The command we need to execute is `gructl apply`.

```bash
$ sudo gructl apply --siterepo site/ memcached
```

Executing the above command generates the following output from our
resources.

```bash
$ sudo gructl apply --siterepo site/ memcached
Loaded 5 resources from 1 modules
package[memcached] is absent, should be present
package[memcached] installing package
package[memcached] resolving dependencies...
package[memcached] looking for conflicting packages...
package[memcached]
package[memcached] Packages (1) memcached-1.4.25-1
package[memcached]
package[memcached] Total Installed Size:  0.14 MiB
package[memcached]
package[memcached] :: Proceed with installation? [Y/n]
package[memcached] checking keyring...
package[memcached] checking package integrity...
package[memcached] loading package files...
package[memcached] checking for file conflicts...
package[memcached] checking available disk space...
package[memcached] :: Processing package changes...
package[memcached] installing memcached...
package[memcached] Optional dependencies for memcached
package[memcached]     perl: for memcached-tool usage [installed]
package[memcached] :: Running post-transaction hooks...
package[memcached] (1/1) Updating manpage index...
package[memcached]
file[/etc/systemd/system/memcached.service.d] is absent, should be present
file[/etc/systemd/system/memcached.service.d] creating resource
file[/etc/systemd/system/memcached.service.d/override.conf] is absent, should be present
file[/etc/systemd/system/memcached.service.d/override.conf] creating resource
shell[systemctl daemon-reload] is absent, should be present
shell[systemctl daemon-reload] executing command
service[memcached] is stopped, should be running
service[memcached] starting service
service[memcached] systemd job id 2013 result: done
```

From the output we can also see that the order of execution of our
resources is correct as we've also seen from the DAG graph in the
previous section.

One last thing we need to do is check the status of the service.

```bash
$ sudo systemctl status memcached
● memcached.service - Memcached Daemon
   Loaded: loaded (/usr/lib/systemd/system/memcached.service; enabled; vendor preset: disabled)
  Drop-In: /etc/systemd/system/memcached.service.d
           └─override.conf
   Active: active (running) since Wed 2016-05-25 17:57:04 EEST; 3min 11s ago
 Main PID: 5613 (memcached)
    Tasks: 6 (limit: 512)
   CGroup: /system.slice/memcached.service
           └─5613 /usr/bin/memcached

May 25 17:57:04 mnikolov-laptop systemd[1]: Started Memcached Daemon.
```

Everything looks good and we can see our drop-in unit being used as well.

## Push configuration

TODO: Add instructions how to start minions
TODO: Add instructions how to push configuration

## View results

TODO: Add instructions how to view results
