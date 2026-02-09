data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

data "snapcd_runner" "myrunner" {
  name = "myrunner"
}

resource "snapcd_module" "mymodule" {
  name                = "mymodule"
  namespace_id        = snapcd_namespace.mynamespace.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module2"
  runner_id           = data.snapcd_runner.myrunner.id
}



resource "snapcd_module_input_from_literal" "myparam" {
  input_kind    = "Param"
  name          = "myvar"
  literal_value = "This will be the value input into 'var.myvar'!"
  module_id     = snapcd_module.mymodule.id
}
