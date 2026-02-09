data "snapcd_runner" "myrunner" {
  name = "myrunner"
}

data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_runner_stack_assignment" "myrunner_mystack" {
  runner_id = data.snapcd_runner.myrunner.id
  stack_id  = data.snapcd_stack.mystack.id
}
