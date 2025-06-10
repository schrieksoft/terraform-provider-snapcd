package secret

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret_store"
	"terraform-provider-snapcd/internal/tests/secret_store_assignment"
)

var AzureKeyVaultSecretScopedToNamespaceCreateConfig = core.NamespaceCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + secret_store_assignment.AzureKeyVaultSecretStoreNamespaceAssignmentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret_scoped_to_namespace" "this" { 
  depends_on         = [snapcd_secret_store_namespace_assignment.this]
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_azure_key_vault_secret_store.this.id
  namespace_id 	     = snapcd_namespace.this.id
}`)

var AzureKeyVaultSecretScopedToModuleCreateConfig = core.ModuleCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + secret_store_assignment.AzureKeyVaultSecretStoreModuleAssignmentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret_scoped_to_module" "this" { 
  depends_on         = [snapcd_secret_store_module_assignment.this]
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_azure_key_vault_secret_store.this.id
  module_id 	       = snapcd_module.this.id
}`)

var AzureKeyVaultSecretScopedToStackCreateConfig = core.StackCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + secret_store_assignment.AzureKeyVaultSecretStoreStackAssignmentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret_scoped_to_stack" "this" { 
  depends_on         = [snapcd_secret_store_stack_assignment.this]
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

var SimpleSecretScopedToNamespaceCreateConfig = core.NamespaceCreateConfig + secret_store.SimpleSecretStoreCreateConfig + secret_store_assignment.SimpleSecretStoreNamespaceAssignmentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_simple_secret_scoped_to_namespace" "this" { 
  depends_on         = [snapcd_secret_store_namespace_assignment.this]
  name  		         = "somevalue%s"
  value              = "somevalue"
  secret_store_id    = snapcd_simple_secret_store.this.id
  namespace_id 	     = snapcd_namespace.this.id
}`)

var SimpleSecretScopedToModuleCreateConfig = core.ModuleCreateConfig + secret_store.SimpleSecretStoreCreateConfig + secret_store_assignment.SimpleSecretStoreModuleAssignmentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_simple_secret_scoped_to_module" "this" { 
  depends_on         = [snapcd_secret_store_module_assignment.this]
  name  		         = "somevalue%s"
  value              = "somevalue"
  secret_store_id    = snapcd_simple_secret_store.this.id
  module_id 	       = snapcd_module.this.id
}`)

var SimpleSecretScopedToStackCreateConfig = core.StackCreateConfig + secret_store.SimpleSecretStoreCreateConfig + secret_store_assignment.SimpleSecretStoreStackAssignmentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_simple_secret_scoped_to_stack" "this" { 
  depends_on         = [snapcd_secret_store_stack_assignment.this]
  name  		         = "somevalue%s"
  value              = "somevalue"
  secret_store_id    = snapcd_simple_secret_store.this.id
  stack_id 	         = snapcd_stack.this.id
}`)

var SimpleSecretCreateConfig = core.StackCreateConfig + secret_store.SimpleSecretStoreCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_simple_secret" "this" {
  name  		         = "somevalue%s"
  value              = "somevalue"
  secret_store_id    = snapcd_simple_secret_store.this.id
}`)



var AwsSecretsManagerSecretScopedToNamespaceCreateConfig = core.NamespaceCreateConfig + secret_store.AwsSecretsManagerSecretStoreCreateConfig + secret_store_assignment.AwsSecretsManagerSecretStoreNamespaceAssignmentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_aws_secrets_manager_secret_scoped_to_namespace" "this" { 
  depends_on         = [snapcd_secret_store_namespace_assignment.this]
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_aws_secrets_manager_secret_store.this.id
  namespace_id 	     = snapcd_namespace.this.id
}`)

var AwsSecretsManagerSecretScopedToModuleCreateConfig = core.ModuleCreateConfig + secret_store.AwsSecretsManagerSecretStoreCreateConfig + secret_store_assignment.AwsSecretsManagerSecretStoreModuleAssignmentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_aws_secrets_manager_secret_scoped_to_module" "this" { 
  depends_on         = [snapcd_secret_store_module_assignment.this]
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_aws_secrets_manager_secret_store.this.id
  module_id 	       = snapcd_module.this.id
}`)

var AwsSecretsManagerSecretScopedToStackCreateConfig = core.StackCreateConfig + secret_store.AwsSecretsManagerSecretStoreCreateConfig + secret_store_assignment.AwsSecretsManagerSecretStoreStackAssignmentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_aws_secrets_manager_secret_scoped_to_stack" "this" { 
  depends_on         = [snapcd_secret_store_stack_assignment.this]
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_aws_secrets_manager_secret_store.this.id
  stack_id 	         = snapcd_stack.this.id
}`)

var AwsSecretsManagerSecretCreateConfig = core.StackCreateConfig + secret_store.AwsSecretsManagerSecretStoreCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_aws_secrets_manager_secret" "this" { 
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_aws_secrets_manager_secret_store.this.id
}`)
