---
page_title: "snapcd_simple_secret_store Resource - snapcd"
subcategory: "Secret Stores"
description: |-
  Manages a Simple Secret Store in Snap CD.
---

# snapcd_simple_secret_store (Resource)

Manages a Simple Secret Store in Snap CD.


## Example Usage

```terraform
resource "snapcd_simple_secret_store" "mysecretstore" {
  name = "mysecretstore"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Unique Name of the Secret Store.

### Optional

- `is_assigned_to_all_scopes` (Boolean) If set to true, secrets scoped to any resource in the system (any Stack, Namespace, Module or Output) can be assigned to this Secret Store

### Read-Only

- `id` (String) Unique ID of the Secret Store.

## Import

Import is supported using the following syntax:

```shell
RESOURCE_ID="12345678-90ab-cdef-1234-56789abcdef0"
terraform import snapcd_simple_secret_store.this $RESOURCE_ID
```
