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


// If you declare a Namespace Param called "myvar_declared_on_namespace", you can map it to "myvar" as follows:
// Note that you can do this mapping irrespective of whether the Namespace Param's "Usage Mode" was set to "UseByDefault" or "UseIfSelected". 
// However, it it was set to "UseByDefault", both "var.myvar_declared_on_namespace" and "var.myvar" will be provided as input variables
// when the Module executes
resource "snapcd_module_param_from_namespace" "myparam" {
  name           = "myvar"
  reference_name = "myvar_declared_on_namespace"
  module_id      = snapcd_module.mymodule.id
}
