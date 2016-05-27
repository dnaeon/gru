## service

The `service` resource is used for managing services on GNU/Linux
systems, which run with systemd.

## Embeds

* [base](base.md)

## Parameters

### name

Name of the service to manage.

* Type: string
* Required: No
* Default: Defaults to the resource title

### enable

Boolean flag indicating whether to enable or disable the service
during boot-time.

* Type: bool
* Required: No
* Default: true

## Example usage

```hcl
service "sshd" {
  state = "running"
  enable = true
}
```
