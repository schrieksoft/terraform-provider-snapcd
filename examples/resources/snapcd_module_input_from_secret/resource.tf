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


// Create a Secret Store and assign it to the "default" Stack. This allows secrets scoped to the "default"
// Stack (or scoped to any of its Namespaces or Modules) to be stored in this Secret Store
resource "snapcd_azure_key_vault_secret_store" "mysecretstore" {
  name          = "my-secret-store"
  key_vault_url = "https://name-of-my-azure-key-vault.vault.azure.net/"
}

resource "secret_store_stack_assignment" "mysecretstore_default" {
  secret_store_id = snapcd_azure_key_vault_secret_store.mysecretstore.id
  stack_id        = data.snapcd_stack.default.id
}

// Create a Secret that references a secret stored in an Azure Key Vault. Scope the Secret to the "default" Stack, 
// meaning all modules within the Stack can access it.
resource "snapcd_azure_key_vault_secret_scoped_to_stack" "mysecret" {
  depends_on         = [secret_store_stack_assignment.mysecretstore_default]
  name               = "my-secret"
  remote_secret_name = "the-name-of-the-secret-in-azure-key-vault" // NOTE this secret must created in the Azure Key Vault separately
  stack_id           = data.snapcd_stack.default.id
  secret_store_id    = snapcd_azure_key_vault_secret_store.mysecretstore.id

}


// Map the contents of "my-secret" to the input variable "var.myvar"
resource "snapcd_module_input_from_secret" "myparam" {
  input_kind = "Param"
  name       = "myvar"
  secret_id  = snapcd_azure_key_vault_secret_scoped_to_stack.mysecret.id
  module_id  = snapcd_module.mymodule.id
}
