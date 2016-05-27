## pacman

The `pacman` resource is used for package management on Arch Linux
systems.

## Embeds

* [base](base.md)

## Parameters

### name

Name of the package.

* Type: string
* Required: yes
* Default: defaults to the resource title

### state

Desired state of the package.

* Type: string
* Required: no
* Default: present

## Example usage

```hcl
pacman "tmux" {
  state = "present"
}
```
