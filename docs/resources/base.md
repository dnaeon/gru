## base

The `base` resource provides a set of base parameters used by all
resources, which embed it.

This resource is not meant to be used directly in your modules, but
instead it is used for providing a common set of parameters shared
between your resources.

## Parameters

### state

Desired state of the resource.

* Type: string
* Required: No
* Default: none
* Values: `present`, `absent`, `running`, `stopped`

### before

An array of resource identifiers before which the current resource
should be evaluated and processed.

* Type: array
* Required: No
* Default: none

### require

An array of resource identifiers after which the current resource
should be evaluated and processed.

* Type: array
* Required: No
* Default: none
