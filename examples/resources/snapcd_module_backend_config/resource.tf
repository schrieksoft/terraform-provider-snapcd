data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

resource "snapcd_runner" "myrunnerpool" {
  name = "myrunnerpool"
}

resource "snapcd_module" "mymodule" {
  name                = "mymodule"
  namespace_id        = snapcd_namespace.mynamespace.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module2"
  runner_id           = snapcd_runner.myrunnerpool.id
}

# The below will result in:
# terraform init \
#   -backend-config="bucket=my-terraform-state-bucket" \
#   -backend-config="key=state/myproject.tfstate" \
#   -backend-config="region=us-east-1"


resource "snapcd_module_backend_config" "bucket" {
  module_id = snapcd_module.mymodule.id
  name      = "bucket"
  value     = "my-terraform-state-bucket"
}

resource "snapcd_module_backend_config" "key" {
  module_id = snapcd_module.mymodule.id
  name      = "key"
  value     = "state/myproject.tfstate"
}

resource "snapcd_module_backend_config" "region" {
  module_id = snapcd_module.mymodule.id
  name      = "region"
  value     = "us-east-1"
}
