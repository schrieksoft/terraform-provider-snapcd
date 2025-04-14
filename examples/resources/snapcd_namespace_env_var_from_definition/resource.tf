data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}



resource "snapcd_namespace_env_var_from_definition" "myenvvar" {
  name            = "NAMESPACE_NAME"
  definition_name = "ModuleName"
  namespace_id    = snapcd_namespace.mynamespace.id
  usage_mode      = "UseByDefault"
}

resource "snapcd_namespace_env_var_from_definition" "myenvvar" {
  name            = "MODULE_NAME"
  definition_name = "ModuleName"
  namespace_id    = snapcd_namespace.mynamespace.id
  usage_mode      = "UseByDefault"
}
