data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}



resource "snapcd_namespace_input_from_definition" "myparam" {
  input_kind      = "Param"
  name            = "myvar"
  definition_name = "ModuleName"
  namespace_id    = snapcd_namespace.mynamespace.id
  usage_mode      = "UseByDefault"
}
