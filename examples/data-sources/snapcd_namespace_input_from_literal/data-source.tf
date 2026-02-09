data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}
data "snapcd_namespace_input_from_literal" "myvar" {
  input_kind   = "Param"
  name         = "myvar"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
