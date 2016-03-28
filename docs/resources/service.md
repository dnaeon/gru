## Service

The `service` resource is used for managing services on a GNU/Linux
system running systemd.

## Parameters

### name

Name of the service resource.

* Type: string
* Required: Yes

### state

Desired state of the service.

* Type: string
* Required: No
* Default: running

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
