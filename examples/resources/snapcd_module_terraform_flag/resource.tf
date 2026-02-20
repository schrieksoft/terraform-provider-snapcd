resource "snapcd_module_terraform_flag" "refresh" {
  module_id = snapcd_module.mymodule.id
  task      = "Plan"
  flag      = "Refresh"
  value     = "false"
}

resource "snapcd_module_terraform_flag" "no_color" {
  module_id = snapcd_module.mymodule.id
  task      = "Plan"
  flag      = "NoColor"
}
