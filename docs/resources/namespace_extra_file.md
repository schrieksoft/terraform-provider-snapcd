---
page_title: "snapcd_namespace_extra_file Resource - snapcd"
subcategory: "Extra Files"
description: |-
  Manages a Namespace Extra File in Snap CD.
---

# snapcd_namespace_extra_file (Resource)

Manages a Namespace Extra File in Snap CD.


## Example Usage

```terraform
data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}


// Add two "Extra Files" to the module. You can add any files you need here. This specific sample solves the following:

// The parameterisation for provider initialization is typically done in the root of a Terraform project and then passed done. 
// As such a pure "module" definition will not have have a "providers.tf" file, nor will it typically have variables with which
// to populate such a file. The below example provides a simple approach for how you could initialize such a module directly with
// Snap CD by passing in Extra Files.

// The following will add the files "provider.tf" and "providers_vars.tf" as "Extra Files", meaning that when any Module in this
// Namespace executes, these two files will be added to the root folder of the Module. 

resource "snapcd_namespace_extra_file" "provider_vars" {
  file_name    = "providers_vars.tf"
  contents     = <<EOT
variable "subscription_id" {}
variable "tenant_id" {}
  EOT
  namespace_id = snapcd_namespace.mynamespace.id
}

resource "snapcd_namespace_extra_file" "providers" {
  file_name    = "providers.tf"
  contents     = <<EOT
provider "azurerm" {
  subscription_id = var.subscription_id
  tenant_id       = var.tenant_id
  features {}
}
  EOT
  namespace_id = snapcd_namespace.mynamespace.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `contents` (String) Contents of the Namespace Extra File
- `file_name` (String) Name of the Namespace Extra File. This name will be use as the name of the file that is created. Must be unique in combination with `namespace_id`.
- `namespace_id` (String) ID of the Namespace Extra File's parent Namespace.

### Optional

- `overwrite` (Boolean) If true any pre-existing file with the same name will be overwritten.

### Read-Only

- `id` (String) Unique ID of the Namespace Extra File.

## Import

Import is supported using the following syntax:

```shell
RESOURCE_ID="12345678-90ab-cdef-1234-56789abcdef0"
terraform import snapcd_namespace_extra_file.this $RESOURCE_ID
```
