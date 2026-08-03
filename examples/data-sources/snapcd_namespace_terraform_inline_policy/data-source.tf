data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}
data "snapcd_namespace_terraform_inline_policy" "mypolicy" {
  name         = "mypolicy"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
