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
  default_engine      = "OpenTofu"

  // example of how to set optional (default) backend args for "init"
  default_init_backend_args = <<EOT
    -backend-config="storage_account_name=somestorageaccount" \
    -backend-config="container_name=terraform-states" \
    -backend-config="key=mystatefile.tfstate" \
    -backend-config="resource_group_name=someresourcegroup" \
    -backend-config="subscription_id=xxxx-xxx-xxx-xxx-xxxx" \
    -backend-config="tenant_id=zzzz-zzz-zzz-zzz-zzzzzz"
  EOT
}


// If you declare a Namespace Param called "myvar_declared_on_namespace", you can map it to "myvar" as follows:
// Note that you can do this mapping irrespective of whether the Namespace Param's "Usage Mode" was set to "UseByDefault" or "UseIfSelected". 
// However, it it was set to "UseByDefault", both "var.myvar_declared_on_namespace" and "var.myvar" will be provided as input variables
// when the Module executes
resource "snapcd_module_input_from_namespace" "myparam" {
  input_kind     = "Param"
  name           = "myvar"
  reference_name = "myvar_declared_on_namespace"
  module_id      = snapcd_module.mymodule.id
}
