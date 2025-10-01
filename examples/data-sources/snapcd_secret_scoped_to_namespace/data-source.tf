data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

data "snapcd_simple_secret_scoped_to_namespace" "mysecret" {
  name         = "my-secret"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
