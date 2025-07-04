---
page_title: "snapcd_namespace_input_from_definition Resource - snapcd"
subcategory: "Namespace Inputs"
description: |-
  Manages a Namespace Input (From Definition) in Snap CD.
---

# snapcd_namespace_input_from_definition (Resource)

Manages a Namespace Input (From Definition) in Snap CD.


## Example Usage

```terraform
data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}



resource "snapcd_namespace_input_from_definition" "myparam" {
  input_kind      = "Param"
  name            = "myvar"
  definition_name = "ModuleName"
  namespace_id    = snapcd_namespace.mynamespace.id
  usage_mode      = "UseByDefault"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `definition_name` (String) Name of the Definition from which to get take the input. Must be one of 'ModuleId', 'NamespaceId', 'StackId', 'ModuleName', 'NamespaceName', 'StackName', 'SourceUrl', 'SourceRevision' and 'SourceSubdirectory'
- `input_kind` (String) The kind of input. Must be one of 'Param' or 'EnvVar'. Changing this will force the resource to be recreated.
- `name` (String) Name of the Namespace Input (From Definition).  Must be unique in combination with `namespace_id`.
- `namespace_id` (String) ID of the Namespace Input (From Definition)'s parent Namespace.

### Optional

- `usage_mode` (String) Whether the input should be used by default on all Modules, or only when explicitly selected on the Module itself. Must be one of 'UseIfSelected' and 'UseByDefault'

### Read-Only

- `id` (String) Unique ID of the Namespace Input (From Definition).

## Import

Import is supported using the following syntax:

```shell
RESOURCE_ID="12345678-90ab-cdef-1234-56789abcdef0"
terraform import snapcd_namespace_param_from_definition.this $RESOURCE_ID
```
