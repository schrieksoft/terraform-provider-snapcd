resource "snapcd_module_hook" "before_plan" {
  module_id = snapcd_module.mymodule.id
  task      = "Plan"
  phase     = "Before"
  script    = "echo 'about to plan'"
}

resource "snapcd_module_hook" "after_apply" {
  module_id = snapcd_module.mymodule.id
  task      = "Apply"
  phase     = "After"
  script    = "echo 'apply finished'"
}
