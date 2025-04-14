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

resource "snapcd_runner_pool_namespace_assignment" "mysp_administrator" {
  runner_pool_id = snapcd_runner_pool.myrunnerpool.id
  namespace_id   = snapcd_namespace.mynamespace.id
}
