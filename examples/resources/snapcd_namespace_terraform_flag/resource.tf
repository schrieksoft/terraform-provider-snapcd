resource "snapcd_namespace_terraform_flag" "refresh" {
  namespace_id = snapcd_namespace.mynamespace.id
  task         = "Plan"
  flag         = "Refresh"
  value        = "false"
}

resource "snapcd_namespace_terraform_flag" "no_color" {
  namespace_id = snapcd_namespace.mynamespace.id
  task         = "Plan"
  flag         = "NoColor"
}
