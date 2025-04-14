data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}



resource "snapcd_namespace_param_from_definition" "myparam" {
  name            = "myvar"
  definition_name = "ModuleName"
  namespace_id    = snapcd_namespace.mynamespace.id
  usage_mode      = "UseByDefault"
}
