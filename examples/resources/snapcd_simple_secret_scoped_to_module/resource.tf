data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

resource "snapcd_runner_pool" "myrunnerpool" {
  name = "myrunnerpool"
}

resource "snapcd_module" "mymodule" {
  name                = "mymodule"
  namespace_id        = snapcd_namespace.mynamespace.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module2"
  runner_pool_id      = data.snapcd_runner_pool.default.id
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


// Create a Secret and scope it to the Secret to the "mymodule" Module, meaning that only this module can access it.
resource "snapcd_simple_secret_scoped_to_module" "mysecret" {
  name            = "my-secret"
  value           = "something-very-secret" // NOTE that creating a secret like this means that this value will be stored in the terraform state file!
  module_id       = snapcd_module.mymodule.id
  secret_store_id = snapcd_simple_secret_store.mysecretstore.id

}
