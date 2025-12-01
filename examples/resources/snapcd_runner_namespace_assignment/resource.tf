data "snapcd_stack" "mystack" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

data "snapcd_runner" "myrunner" {
  name = "myrunner"
}

resource "snapcd_runner_namespace_assignment" "myrunner_mynamespace" {
  runner_id    = data.snapcd_runner.myrunner.id
  namespace_id = snapcd_namespace.mynamespace.id
}
