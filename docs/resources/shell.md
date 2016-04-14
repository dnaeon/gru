## shell

The `shell` resource is used for executing shell commands.

The command that is to be executed should be idempotent. If the
command that is to be executed is not idempotent on it's own,
in order to achieve idempotency of the resource you should set the
`creates` parameter to a filename that can be checked for existence.

## Parameters

### command

Command to be executed.

* Type: string
* Required: No
* Default: Defaults to the resource title

### creates

A filename to check for existence before running the command.

You should use `creates` parameter to achieve idempotency of your
resource, if the command to be executed is not idempotent on it's own.

* Type: string
* Required: No

## Example usage

```hcl
shell "touch /tmp/foo" {
  creates = "/tmp/foo"
}
```
