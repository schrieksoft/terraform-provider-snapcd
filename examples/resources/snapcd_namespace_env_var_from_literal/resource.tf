data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}


resource "snapcd_namespace_env_var_from_literal" "myenvvar" {
  name          = "SOME_ENV_VAR"
  literal_value = "This will be the value of the env var!"
  namespace_id  = snapcd_namespace.mynamespace.id
  usage_mode    = "UseByDefault"
}
