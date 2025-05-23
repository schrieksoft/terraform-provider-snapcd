---
page_title: "snapcd_module_param_from_literal Data Source - snapcd"
subcategory: "Module Inputs (Parameters)"
description: |-
  Use this data source to access information about an existing Module Param (From Literal) in Snap CD.
---

# snapcd_module_param_from_literal (Data Source)

Use this data source to access information about an existing Module Param (From Literal) in Snap CD.


## Example Usage

```terraform
data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}
data "snapcd_module" "mymodule" {
  name         = "mymodule"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
data "snapcd_module_param_from_literal" "myvar" {
  name      = "myvar"
  module_id = data.snapcd_module.mymodule.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `module_id` (String) ID of the Module Param (From Literal)'s parent Module.
- `name` (String) Name of the Module Param (From Literal).  Must be unique in combination with `module_id`.

### Read-Only

- `id` (String) Unique ID of the Module Param (From Literal).
- `literal_value` (String) Literal value of the input.
- `type` (String) Type of literal input. Must be one of 'String' and 'NotString'. Use 'NotString' for values such as numbers, bools, list, maps etc.
