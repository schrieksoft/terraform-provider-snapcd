data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
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
