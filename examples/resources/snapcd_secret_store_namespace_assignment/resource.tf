data "snapcd_stack" "default" {
  name = "default"
}

resource "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}

resource "snapcd_azure_key_vault_secret_store" "mysecretstore" {
  name          = "mysecretstore"
  key_vault_url = "https://name-of-my-azure-key-vault.vault.azure.net/"
}

// Assign Secret Store to Namespace, allowing to create Secrets scoped to this Namespace (or any of its child Modules) to be created in this Secret Store
resource "snapcd_secret_store_namespace_assignment" "mysecretstore_mynamespace" {
  secret_store_id = snapcd_azure_key_vault_secret_store.mysecretstore.id
  namespace_id    = snapcd_namespace.mynamespace.id
}
