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

// Provided you have a Module called "anothermodule" within a namespace called "anothernamespace" (within the same Stack as "mymodule"),
// which provides an output called "some_output", you can set the input parameter "var.myvar" equal to the value stored in "some_output" as
// follows:
data "snapcd_namespace" "anothernamespace" {
  name     = "anothernamespace"
  stack_id = data.snapcd_stack.mystack.id
}

data "snapcd_module" "anothermodule" {
  name         = "anothermodule"
  namespace_id = data.snapcd_namespace.anothernamespace.id
}

resource "snapcd_module_input_from_output" "myparam" {
  input_kind       = "Param"
  name             = "myvar"
  output_name      = "some_output"
  output_module_id = data.snapcd_module.anothermodule.id
  module_id        = snapcd_module.mymodule.id
}
