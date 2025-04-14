data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_stack" "mynamespace" {
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
}


// In the below example, the name of the module (i.e. the value in snapcd_module.mymodule.name) will be bound to the "SNAPCD_ENV_MODULE_NAME" environment variable when this module executes on the Runner.
resource "snapcd_module_env_var_from_definition" "myenvvar" {
  name            = "MODULE_NAME"
  definition_name = "ModuleName"
  module_id       = snapcd_namespace.mymodule.id
}
