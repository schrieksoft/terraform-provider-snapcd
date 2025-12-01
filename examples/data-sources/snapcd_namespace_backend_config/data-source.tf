data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

# Namespace Backend Config Data Source Example
# This example retrieves information about an existing
# namespace backend configuration in SnapCD

data "snapcd_namespace_backend_config" "bucket" {
  namespace_id = data.snapcd_namespace.mynamespace.id
  name         = "bucket"
}

data "snapcd_namespace_backend_config" "key" {
  namespace_id = data.snapcd_namespace.mynamespace.id
  name         = "key"
}

data "snapcd_namespace_backend_config" "region" {
  namespace_id = data.snapcd_namespace.mynamespace.id
  name         = "region"
}

# Use the retrieved backend config values in outputs
output "backend_bucket" {
  value = data.snapcd_namespace_backend_config.bucket.value
}

output "backend_key" {
  value = data.snapcd_namespace_backend_config.key.value
}

output "backend_region" {
  value = data.snapcd_namespace_backend_config.region.value
}
