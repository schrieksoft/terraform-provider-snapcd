resource "snapcd_module_pulumi_flag" "refresh" {
  module_id = snapcd_module.mymodule.id
  task      = "Plan"
  flag      = "Refresh"
  value     = "true"
}

resource "snapcd_module_pulumi_flag" "debug" {
  module_id = snapcd_module.mymodule.id
  task      = "Plan"
  flag      = "Debug"
}
