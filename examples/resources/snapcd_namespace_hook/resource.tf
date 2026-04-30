resource "snapcd_namespace_hook" "before_plan" {
  namespace_id = snapcd_namespace.mynamespace.id
  task         = "Plan"
  phase        = "Before"
  script       = "echo 'namespace-level before plan'"
}
