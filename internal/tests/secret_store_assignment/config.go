package secret_store_assignment

var AzureKeyVaultSecretStoreStackAssignmentCreateConfig = `
resource "snapcd_secret_store_stack_assignment" "this" { 
  secret_store_id = snapcd_azure_key_vault_secret_store.this.id
  stack_id        = snapcd_stack.this.id
}`

var AzureKeyVaultSecretStoreNamespaceAssignmentCreateConfig = `
resource "snapcd_secret_store_namespace_assignment" "this" { 
  secret_store_id = snapcd_azure_key_vault_secret_store.this.id
  namespace_id    = snapcd_namespace.this.id
}`

var AzureKeyVaultSecretStoreModuleAssignmentCreateConfig = `
resource "snapcd_secret_store_module_assignment" "this" { 
  secret_store_id = snapcd_azure_key_vault_secret_store.this.id
  module_id        = snapcd_module.this.id
}`

var SimpleSecretStoreStackAssignmentCreateConfig = `
resource "snapcd_secret_store_stack_assignment" "this" { 
  secret_store_id = snapcd_simple_secret_store.this.id
  stack_id        = snapcd_stack.this.id
}`

var SimpleSecretStoreNamespaceAssignmentCreateConfig = `
resource "snapcd_secret_store_namespace_assignment" "this" { 
  secret_store_id = snapcd_simple_secret_store.this.id
  namespace_id    = snapcd_namespace.this.id
}`

var SimpleSecretStoreModuleAssignmentCreateConfig = `
resource "snapcd_secret_store_module_assignment" "this" { 
  secret_store_id = snapcd_simple_secret_store.this.id
  module_id        = snapcd_module.this.id
}`

var AwsSecretsManagerSecretStoreStackAssignmentCreateConfig = `
resource "snapcd_secret_store_stack_assignment" "this" { 
  secret_store_id = snapcd_aws_Secrets_Manager_secret_store.this.id
  stack_id        = snapcd_stack.this.id
}`

var AwsSecretsManagerSecretStoreNamespaceAssignmentCreateConfig = `
resource "snapcd_secret_store_namespace_assignment" "this" { 
  secret_store_id = snapcd_aws_Secrets_Manager_secret_store.this.id
  namespace_id    = snapcd_namespace.this.id
}`

var AwsSecretsManagerSecretStoreModuleAssignmentCreateConfig = `
resource "snapcd_secret_store_module_assignment" "this" { 
  secret_store_id = snapcd_aws_Secrets_Manager_secret_store.this.id
  module_id        = snapcd_module.this.id
}`
