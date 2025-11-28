data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

resource "snapcd_runner" "myrunnerpool" {
  name = "myrunnerpool"
}

resource "snapcd_runner_namespace_assignment" "myrunnerpool_mynamespace" {
  runner_id    = snapcd_runner.myrunnerpool.id
  namespace_id = snapcd_namespace.mynamespace.id
}
