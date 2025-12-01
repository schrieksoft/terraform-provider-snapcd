data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

data "snapcd_module" "mymodule" {
  name         = "mymodule"
  namespace_id = data.snapcd_namespace.mynamespace.id
}

# Module Backend Config Data Source Example
# This example retrieves information about an existing
# module backend configuration in SnapCD

data "snapcd_module_backend_config" "bucket" {
  module_id = data.snapcd_module.mymodule.id
  name      = "bucket"
}

data "snapcd_module_backend_config" "key" {
  module_id = data.snapcd_module.mymodule.id
  name      = "key"
}

data "snapcd_module_backend_config" "region" {
  module_id = data.snapcd_module.mymodule.id
  name      = "region"
}

# Use the retrieved backend config values in outputs
output "backend_bucket" {
  value = data.snapcd_module_backend_config.bucket.value
}

output "backend_key" {
  value = data.snapcd_module_backend_config.key.value
}

output "backend_region" {
  value = data.snapcd_module_backend_config.region.value
}
