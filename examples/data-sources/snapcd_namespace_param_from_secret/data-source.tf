data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}
data "snapcd_namespace_param_from_secret" "myvar" {
  name         = "myvar"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
