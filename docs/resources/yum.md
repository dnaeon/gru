## yum

The `yum` resource is used for package management under
RHEL/CentOS systems.

## Embeds

* [base](base.md)

## Parameters

### name

Name of the package.

* Type: string
* Required: yes
* Default: defaults to the resource title

### version

Version of the package.

* Type: string
* Required: no
* Default: none

## Example usage

```hcl
yum "tmux" {
  state = "present"
}
```
