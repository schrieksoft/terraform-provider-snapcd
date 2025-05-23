---
page_title: "snapcd_resource_role_assignment Resource - snapcd"
subcategory: "Identity Access Management"
description: |-
  Manages a Resource Role Assignment in Snap CD.
---

# snapcd_resource_role_assignment (Resource)

Manages a Resource Role Assignment in Snap CD.


## Example Usage

```terraform
data "snapcd_service_principal" "mysp" {
  client_id = "MyServicePrincipal"
}

data "snapcd_stack" "default" {
  name = "default"
}


resource "snapcd_resource_role_assignment" "mysp_administrator" {
  principal_id            = snapcd_service_principal.mysp.id
  resource_id             = data.snapcd_stack.default.id
  principal_discriminator = "ServicePrincipal"
  resource_discriminator  = "Stack"
  role_name               = "Owner"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `principal_discriminator` (String) Type of Principal that the `principal_id` identifies. Must be one of 'User', 'ServicePrincipal' and 'Group'
- `principal_id` (String) ID of the Principal to which the role is assigned.
- `resource_discriminator` (String) Type of Resource that the `resource_id` identifies.
- `resource_id` (String) ID of the Resource on which the role applies.
- `role_name` (String) Name of the Role that is assigned.Must be one of 'Owner', 'Contributor' and 'Reader'

### Read-Only

- `id` (String) Unique ID of the Resource Role Assignment.

## Import

Import is supported using the following syntax:

```shell
RESOURCE_ID="12345678-90ab-cdef-1234-56789abcdef0"
terraform import snapcd_resource_role_assignment.this $RESOURCE_ID
```
