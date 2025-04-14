data "snapcd_stack" "default" {
  name = "default"
}

// Create a Secret Store
resource "snapcd_simple_secret_store" "mysecretstore" {
  name = "my-secret-store"

}

// Assign Secret Store to the "default" Stack. This allows secrets scoped to the "default"
// Stack (or scoped to any of its Namespaces or Modules) to be stored in this Secret Store
resource "secret_store_stack_assignment" "mysecretstore_default" {
  secret_store_id = snapcd_simple_secret_store.mysecretstore.id
  stack_id        = data.snapcd_stack.default.id
}

// Create a secret an scope the Secret to the "default" Stack, meaning all Modules within the Stack can access it.
resource "snapcd_simple_secret_scoped_to_stack" "mysecret" {
  name            = "my-secret"
  value           = "something-very-secret" // NOTE that creating a secret like this means that this value will be stored in the terraform state file!
  stack_id        = data.snapcd_stack.default.id
  secret_store_id = snapcd_simple_secret_store.mysecretstore.id
}
