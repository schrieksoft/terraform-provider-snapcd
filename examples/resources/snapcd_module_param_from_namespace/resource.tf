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


// If you declare a Namespace Env Var called "SOME_ENV_VAR_DECLARED_ON_NAMESPACE", you can map it to "SOME_ENV_VAR" (which will show up as "SNAPCD_ENV_SOME_ENV_VAR" on the Runner) as follows:
// Note that you can do this mapping irrespective of whether the Namespace Env Var's "Usage Mode" was set to "UseByDefault" or "UseIfSelected". However, it it was set to "UseByDefault", both
// "SNAPCD_ENV_SOME_ENV_VAR_DECLARED_ON_NAMESPACE" and "SNAPCD_ENV_SOME_ENV_VAR" will be mapped to Env Vars on the runner.
resource "snapcd_module_env_var_from_namespace" "myenvvar" {
  name           = "SOME_ENV_VAR"
  reference_name = "SOME_ENV_VAR_DECLARED_ON_NAMESPACE"
  module_id      = snapcd_module.mymodule.id
}
