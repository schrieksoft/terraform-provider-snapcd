package secret

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret_store"
)

var AzureKeyVaultSecretScopedToNamespaceCreateConfig = core.NamespaceCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret_scoped_to_namespace" "this" { 
  name  		     = "somevalue%s"
  remote_secret_name = "name-in-remote-%s"
  secret_store_id    = snapcd_azure_key_vault_secret_store.this.id
  namespace_id 	     = snapcd_namespace.this.id
}`)
