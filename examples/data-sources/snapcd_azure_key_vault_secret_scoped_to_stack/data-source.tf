
data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_simple_secret_store" "mysecretstore" {
  name = "my-secret-store"
}
data "snapcd_simple_secret_scoped_to_stack" "mysecret" {
  name            = "my-secret"
  stack_id        = data.snapcd_stack.default.id
  secret_store_id = data.snapcd_simple_secret_store.mysecretstore.id
}
