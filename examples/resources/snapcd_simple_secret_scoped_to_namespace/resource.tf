data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
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


// Create a Secret and scope it to to the "mynamespace" Namespace, meaning all Modules within the Namespace can access it.
resource "snapcd_simple_secret_scoped_to_namespace" "mysecret" {
  name            = "my-secret"
  value           = "something-very-secret" // NOTE that creating a secret like this means that this value will be stored in the terraform state file!
  namespace_id    = snapcd_namespace.mynamespace.id
  secret_store_id = snapcd_simple_secret_store.mysecretstore.id

}
