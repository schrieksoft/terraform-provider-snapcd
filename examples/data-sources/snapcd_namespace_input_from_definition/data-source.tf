data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}
data "snapcd_namespace_input_from_definition" "myvar" {
  input_kind   = "Param"
  name         = "myvar"
  namespace_id = snapcd_namespace.mynamespace.id
}
