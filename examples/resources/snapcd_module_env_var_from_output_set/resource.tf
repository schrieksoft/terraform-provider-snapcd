data "snapcd_stack" "default" {
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
}


// Provided you have a Module called "anothermodule" within a namespace called "anothernamespace" (within the same Stack as "mymodule"), 
// you can pull in all outputs from that Module as environment variables on "mymodule"

// If "anothermodule" has outputs called "some_output" and "another_output", then the Runner will set the environment variables
// "SNAPCD_ENV_some_output" and "SNAPCD_ENV_another_output" to the values of the respective outputs.

// NOTE that below we set a "name" equal to "from_anothermodule". For other inputs (from definition, from literal, from secret etc.) the "name"
// field acts both as unique identifier (along with module_id) for the database entity, as well as for the name the environement variable
// takes when used by the runner. In the case of a "from output set" Env Var, this name only acts as unique identifier and plays no further role.

resource "snapcd_module_param_from_output_set" "myenvvar" {
  name           = "from_anothermodule"
  module_name    = "anothermodule"
  namespace_name = "anothernamespace"
  module_id      = snapcd_module.mymodule.id
}
