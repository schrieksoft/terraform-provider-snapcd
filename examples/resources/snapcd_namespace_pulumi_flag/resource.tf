resource "snapcd_namespace_pulumi_flag" "refresh" {
  namespace_id = snapcd_namespace.mynamespace.id
  task         = "Plan"
  flag         = "Refresh"
  value        = "true"
}

resource "snapcd_namespace_pulumi_flag" "debug" {
  namespace_id = snapcd_namespace.mynamespace.id
  task         = "Plan"
  flag         = "Debug"
}
