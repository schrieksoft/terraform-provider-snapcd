data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

data "snapcd_namespace_secret" "mysecret" {
  name         = "my-secret"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
