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

resource "snapcd_azure_key_vault_secret_store" "mysecretstore" {
  name          = "mysecretstore"
  key_vault_url = "https://name-of-my-azure-key-vault.vault.azure.net/"
}


resource "snapcd_module" "mymodule" {
  name                = "mymodule"
  namespace_id        = snapcd_namespace.mynamespace.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module2"
  runner_pool_id      = data.snapcd_runner_pool.default.id
}


// Assign Secret Store to Module, allowing to create Secrets scoped to this Module to be created in this Secret Store
resource "snapcd_secret_store_module_assignment" "mysecretstore_mymodule" {
  secret_store_id = snapcd_azure_key_vault_secret_store.mysecretstore.id
  module_id       = snapcd_module.mymodule.id
}
