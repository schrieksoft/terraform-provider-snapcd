resource "snapcd_azure_key_vault_secret_store" "mysecretstore" {
  name          = "mysecretstore"
  key_vault_url = "https://name-of-my-azure-key-vault.vault.azure.net/"
}
