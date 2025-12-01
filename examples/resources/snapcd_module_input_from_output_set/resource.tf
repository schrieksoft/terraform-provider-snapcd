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
  runner_id           = data.snapcd_runner.default.id
}


// Provided you have a Module called "anothermodule" within a namespace called "anothernamespace" (within the same Stack as "mymodule"), 
// you can pull in all outputs from that Module as input Parameters on "mymodule".

// If "anothermodule" has outputs called "some_output" and "another_output", then the Runner will set the Params "var.some_output"
// and "var.another_output" to the values of the respective outputs.

// NOTE that below we set a "name" equal to "from_anothermodule". For other inputs (from Literal, from Output, from Secret etc.) the "name"
// field acts both as unique identifier (along with module_id) for the database entity, as well as determining the name the Param takes when used
// by the runner. In the case of a "from Output Set" Param, this name only acts as unique identifier and plays no further role.

resource "snapcd_module_input_from_output_set" "myparam" {
  input_kind     = "Param"
  name           = "from_anothermodule"
  module_name    = "anothermodule"
  namespace_name = "anothernamespace"
  module_id      = snapcd_module.mymodule.id
}
