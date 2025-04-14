data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_azure_key_vault_secret_store" "mysecretstore" {
  name          = "mysecretstore"
  key_vault_url = "https://name-of-my-azure-key-vault.vault.azure.net/"
}

// Assign Secret Store to Stack, allowing to create Secrets scoped to this Stack (or any of its child Namespaces or Modules) to be created in this Secret Store
resource "snapcd_secret_store_stack_assignment" "mysp_administrator" {
  secret_store_id = snapcd_azure_key_vault_secret_store.mysecretstore.id
  stack_id        = data.snapcd_stack.default.id
}
