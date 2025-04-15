data "snapcd_stack" "default" {
  name = "default"
}
data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}
data "snapcd_module" "mymodule" {
  name         = "mymodule"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
data "snapcd_simple_secret_store" "mysecretstore" {
  name = "my-secret-store"
}
data "snapcd_simple_secret_scoped_to_module" "mysecret" {
  name            = "my-secret"
  module_id       = data.snapcd_module.mymodule.id
  secret_store_id = data.snapcd_simple_secret_store.mysecretstore.id
}
