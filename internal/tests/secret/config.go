package secret

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret_store"
)

var AzureKeyVaultSecretScopedToNamespaceCreateConfig = core.NamespaceCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret_scoped_to_namespace" "this" { 
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_azure_key_vault_secret_store.this.id
  namespace_id 	     = snapcd_namespace.this.id
}`)



var AzureKeyVaultSecretScopedToModuleCreateConfig = core.ModuleCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret_scoped_to_module" "this" { 
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_azure_key_vault_secret_store.this.id
  module_id 	       = snapcd_module.this.id
}`)

var AzureKeyVaultSecretScopedToStackCreateConfig = core.StackCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret_scoped_to_stack" "this" { 
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_azure_key_vault_secret_store.this.id
  stack_id 	         = snapcd_stack.this.id
}`)


var AzureKeyVaultSecretCreateConfig = core.StackCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret" "this" { 
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_azure_key_vault_secret_store.this.id
}`)
