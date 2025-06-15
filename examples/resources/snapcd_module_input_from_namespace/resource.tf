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


// If you declare a Namespace Input called "SOME_PARAM_DECLARED_ON_NAMESPACE", you can map it to "var.SOME_PARAM" when this Module executes on the runner as follows:
// Note that you can do this mapping irrespective of whether the Namespace Input's "Usage Mode" was set to "UseByDefault" or "UseIfSelected". However, if it was set to "UseByDefault", both
// "var.SOME_PARAM_DECLARED_ON_NAMESPACE" and "var.SOME_PARAM" will be provided as parameters to the runner.
resource "snapcd_module_input_from_namespace" "myparam" {
  input_kind     = "Param"
  name           = "SOME_PARAM"
  reference_name = "SOME_PARAM_DECLARED_ON_NAMESPACE"
  module_id      = snapcd_module.mymodule.id
}
