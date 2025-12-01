data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name                = "mynamespace"
  stack_id            = data.snapcd_stack.mystack.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/simple-namespace"
}

# The below will result in:
# terraform init \
#   -backend-config="bucket=my-terraform-state-bucket" \
#   -backend-config="key=state/myproject.tfstate" \
#   -backend-config="region=us-east-1"

resource "snapcd_namespace_backend_config" "bucket" {
  namespace_id = snapcd_namespace.mynamespace.id
  name         = "bucket"
  value        = "my-terraform-state-bucket"
}

resource "snapcd_namespace_backend_config" "key" {
  namespace_id = snapcd_namespace.mynamespace.id
  name         = "key"
  value        = "state/myproject.tfstate"
}

resource "snapcd_namespace_backend_config" "region" {
  namespace_id = snapcd_namespace.mynamespace.id
  name         = "region"
  value        = "us-east-1"
}
