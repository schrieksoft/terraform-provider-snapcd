package secret_store

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
)

var AzureKeyVaultSecretStoreCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret_store" "this" { 
  name  = "somevalue%s"
  key_vault_url = "https://snapcdlocaltesting.vault.azure.net/"
}`)
