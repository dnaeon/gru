## package

The `package` resourse is a meta resource for package management.

This resource makes it possible to write resources without having to
depend on the underlying package manager for your system, by trying to
determine the most appropriate package provider for you.

## Embeds

* [base](base.md)

## Parameters

### name

Name of the package to manage.

* Type: string
* Required: no
* Default: defaults to the resource title

### version

Version of the package.

* Type: string
* Required: no
* Default: none

### provider

Package provider to use for this resource.

* Type: string
* Required: no
* Default: none
* Values: `pacman`

## Example usage

```hcl
package "tmux" {
  state = "present"
}
```
