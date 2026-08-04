data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

resource "snapcd_namespace_terraform_local_policy" "mypolicy" {
  name         = "mypolicy"
  namespace_id = snapcd_namespace.mynamespace.id

  path = "/etc/snapcd/policies/terraform-baseline"
}
