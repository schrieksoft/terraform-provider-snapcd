---
page_title: "snapcd_namespace_input_from_secret Data Source - snapcd"
subcategory: "Namespace Inputs"
description: |-
  Retrieves a Namespace Input (From Secret) from Snap CD.
---

# snapcd_namespace_input_from_secret (Data Source)

Retrieves a Namespace Input (From Secret) from Snap CD.


## Example Usage

```terraform
data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}
data "snapcd_namespace_input_from_secret" "myvar" {
  input_kind   = "Param"
  name         = "myvar"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `input_kind` (String) The kind of input. Must be one of 'Param' or 'EnvVar'. Changing this will force the resource to be recreated.
- `name` (String) Name of the Namespace Input (From Secret).  Must be unique in combination with `namespace_id`.

### Read-Only

- `id` (String) Unique ID of the Namespace Input (From Secret).
- `namespace_id` (String) ID of the Namespace Input (From Secret)'s parent Namespace.
- `secret_id` (String) ID of the Secret to take as input.
- `type` (String) Type of literal input the secret value should be formatted as. Must be one of 'String' and 'NotString'. Use 'NotString' for values such as numbers, bools, list, maps etc.
- `usage_mode` (String) Whether the input should be used by default on all Modules, or only when explicitly selected on the Module itself. Must be one of 'UseIfSelected' and 'UseByDefault'
