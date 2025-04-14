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
}

// Provided you have a Module called "anothermodule" within a namespace called "anothernamespace" (within the same Stack as "mymodule"), 
// which provides an output valled "some_output", you can set the environment variable "SNAPCD_ENV_SOME_ENV_VAR" equal to the value stored
// in "some_output" as follows:
resource "snapcd_module_param_from_output" "myenvvar" {
  name           = "SOME_ENV_VAR"
  output_name    = "some_output"
  module_name    = "anothermodule"
  namespace_name = "anothernamespace"
  module_id      = snapcd_module.mymodule.id
}
