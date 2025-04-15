data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_runner_pool" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}


resource "snapcd_module" "mymodule" {
  name                = "mymodule"
  namespace_id        = snapcd_namespace.mynamespace.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module2"
  runner_pool_id      = data.snapcd_runner_pool.default.id

  // example of how to set optional backend args for "init"
  init_backend_args = <<EOT
    -backend-config="storage_account_name=somestorageaccount" \
    -backend-config="container_name=terraform-states" \
    -backend-config="key=mystatefile.tfstate" \
    -backend-config="resource_group_name=someresourcegroup" \
    -backend-config="subscription_id=xxxx-xxx-xxx-xxx-xxxx" \
    -backend-config="tenant_id=zzzz-zzz-zzz-zzz-zzzzzz"
  EOT

}
