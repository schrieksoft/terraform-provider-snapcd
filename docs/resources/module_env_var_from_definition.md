---
page_title: "snapcd_module_env_var_from_definition Resource - snapcd"
subcategory: "Module Inputs (Env Vars)"
description: |-
  Manages a Module Env Var (From Definition) in Snap CD.
---

# snapcd_module_env_var_from_definition (Resource)

Manages a Module Env Var (From Definition) in Snap CD.


## Example Usage

```terraform
data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

resource "snapcd_runner_pool" "myrunnerpool" {
  name = "myrunnerpool"
}

resource "snapcd_module" "mymodule" {
  name                = "mymodule"
  namespace_id        = snapcd_namespace.mynamespace.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module2"
  runner_pool_id      = data.snapcd_runner_pool.default.id
}


// In the below example, the name of the module (i.e. the value in snapcd_module.mymodule.name) will be bound to the "SNAPCD_ENV_MODULE_NAME" environment variable when this module executes on the Runner.
resource "snapcd_module_env_var_from_definition" "myenvvar" {
  name            = "MODULE_NAME"
  definition_name = "ModuleName"
  module_id       = snapcd_namespace.mymodule.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `definition_name` (String) Name of the Definition from which to get take the input. Must be one of 'ModuleId', 'NamespaceId', 'StackId', 'ModuleName', 'NamespaceName', 'StackName', 'SourceUrl', 'SourceRevision' and 'SourceSubdirectory'.
- `module_id` (String) ID of the Module Env Var (From Definition)'s parent Module.
- `name` (String) Name of the Module Env Var (From Definition).  Must be unique in combination with `module_id`.

### Read-Only

- `id` (String) Unique ID of the Module Env Var (From Definition).

## Import

Import is supported using the following syntax:

```shell
RESOURCE_ID="12345678-90ab-cdef-1234-56789abcdef0"
terraform import snapcd_module_env_var_from_definition.this $RESOURCE_ID
```
