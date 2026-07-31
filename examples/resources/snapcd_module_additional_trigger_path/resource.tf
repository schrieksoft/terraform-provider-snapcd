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
  source_subdirectory = "modules/app-a"
  runner_id           = data.snapcd_runner.myrunner.id

  trigger_on_source_changed = true
}

// When path-scoped triggering is enabled for the module, it only triggers when a commit touches
// its own source_subdirectory (or a directory it references). Additional Trigger Paths widen that
// watch set with directories the module depends on in ways Snap CD cannot discover statically —
// for example a folder read through a symlink, or inputs consumed by a script.

resource "snapcd_module_additional_trigger_path" "shared_scripts" {
  path      = "shared/scripts"
  module_id = snapcd_module.mymodule.id
}
