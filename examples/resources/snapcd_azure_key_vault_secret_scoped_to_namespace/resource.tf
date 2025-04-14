data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

// Create a Secret Store
resource "snapcd_azure_key_vault_secret_store" "mysecretstore" {
  name          = "my-secret-store"
  key_vault_url = "https://name-of-my-azure-key-vault.vault.azure.net/"

}

// Assign Secret Store to the "default" Stack. This allows secrets scoped to the "default"
// Stack (or scoped to any of its Namespaces or Modules) to be stored in this Secret Store
resource "secret_store_stack_assignment" "mysecretstore_default" {
  secret_store_id = snapcd_azure_key_vault_secret_store.mysecretstore.id
  stack_id        = data.snapcd_stack.default.id
}


// Create a Secret that references a secret stored in an Azure Key Vault. Scope the Secret to the "mynamespace" Namespace, 
// meaning all Modules within the Namespace can access it.
resource "snapcd_azure_key_vault_secret_scoped_to_namespace" "mysecret" {
  name               = "my-secret"
  remote_secret_name = "the-name-of-the-secret-in-azure-key-vault" // NOTE this secret must created in the Azure Key Vault separately
  namespace_id       = snapcd_namespace.mynamespace.id
  secret_store_id    = snapcd_azure_key_vault_secret_store.mysecretstore.id

}
