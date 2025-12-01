data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_stack_secret" "mysecret" {
  name     = "my-secret"
  stack_id = data.snapcd_stack.mystack.id
}
