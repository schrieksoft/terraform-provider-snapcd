data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}
data "snapcd_namespace_env_var_from_secret" "myenvvar" {
  name         = "MY_ENV_VAR"
  namespace_id = snapcd_namespace.mynamespace.id
}
