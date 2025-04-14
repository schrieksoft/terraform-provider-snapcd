data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_runner_pool" "myrunnerpool" {
  name = "myrunnerpool"
}

resource "snapcd_runner_pool_stack_assignment" "mysp_administrator" {
  runner_pool_id = snapcd_runner_pool.myrunnerpool.id
  stack_id       = data.snapcd_stack.default.id
}
