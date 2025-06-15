data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}


// Create a Secret Store and assign it to the "default" Stack. This allows secrets scoped to the "default"
// Stack (or scoped to any of its Namespaces or Modules) to be stored in this Secret Store
resource "snapcd_azure_key_vault_secret_store" "mysecretstore" {
  name          = "my-secret-store"
  key_vault_url = "https://name-of-my-azure-key-vault.vault.azure.net/"

}

resource "snapcd_secret_store_stack_assignment" "mysecretstore_default" {
  secret_store_id = snapcd_azure_key_vault_secret_store.mysecretstore.id
  stack_id        = data.snapcd_stack.default.id
}


// Create a Secret that references a secret stored in an Azure Key Vault.
resource "snapcd_azure_key_vault_secret" "mysecret" {
  depends_on         = [snapcd_secret_store_stack_assignment.mysecretstore_default]
  name               = "my-secret"
  remote_secret_name = "the-name-of-the-secret-in-azure-key-vault" // NOTE this secret must created in the Azure Key Vault separately
  secret_store_id    = snapcd_azure_key_vault_secret_store.mysecretstore.id

}


// Map the contents of "my-secret" to the input parameter var.myvar
resource "snapcd_namespace_param_from_secret" "myparam" {
  name         = "myvar"
  secret_id    = snapcd_azure_key_vault_secret.mysecret.id
  namespace_id = snapcd_namespace.mynamespace.id
  usage_mode   = "UseByDefault"
}
