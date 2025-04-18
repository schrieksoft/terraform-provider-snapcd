---
page_title: "snapcd_simple_secret_scoped_to_namespace Resource - snapcd"
subcategory: "Secrets"
description: |-
  Manages a Simple Secret (Scoped to Namespace) in Snap CD.
---

# snapcd_simple_secret_scoped_to_namespace (Resource)

Manages a Simple Secret (Scoped to Namespace) in Snap CD.


## Example Usage

```terraform
data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

// Create a Secret Store
resource "snapcd_simple_secret_store" "mysecretstore" {
  name = "my-secret-store"

}

// Assign Secret Store to the "default" Stack. This allows secrets scoped to the "default"
// Stack (or scoped to any of its Namespaces or Modules) to be stored in this Secret Store
resource "secret_store_stack_assignment" "mysecretstore_default" {
  secret_store_id = snapcd_simple_secret_store.mysecretstore.id
  stack_id        = data.snapcd_stack.default.id
}


// Create a Secret and scope it to to the "mynamespace" Namespace, meaning all Modules within the Namespace can access it.
resource "snapcd_simple_secret_scoped_to_namespace" "mysecret" {
  name            = "my-secret"
  value           = "something-very-secret" // NOTE that creating a secret like this means that this value will be stored in the terraform state file!
  namespace_id    = snapcd_namespace.mynamespace.id
  secret_store_id = snapcd_simple_secret_store.mysecretstore.id

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Unique Name within of the Secret within the Secret Store.
- `namespace_id` (String) Id of the Namespace to scope the Secret to
- `secret_store_id` (String) Unique ID of the Secret.
- `value` (String, Sensitive) Value of the to store in the Simple Secret Store. NOTE that this value **will** end up in the .tfstate file. If you wish to avoid this, create the secret directly via the API or Dashboard instead.

### Read-Only

- `id` (String) Unique ID of the Secret.

## Import

Import is supported using the following syntax:

```shell
RESOURCE_ID="12345678-90ab-cdef-1234-56789abcdef0"
terraform import snapcd_simple_secret_scoped_to_namespace.this $RESOURCE_ID
```
