data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_stack_secret" "mysecret" {
  name     = "my-secret"
  stack_id = data.snapcd_stack.default.id
}
