data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_runner" "myrunner" {
  name = "myrunner"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

resource "snapcd_module" "mymodule" {
  name                = "mymodule"
  namespace_id        = snapcd_namespace.mynamespace.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module2"
  runner_id           = data.snapcd_runner.myrunner.id
  engine              = "OpenTofu"
}
