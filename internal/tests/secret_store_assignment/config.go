package secret_store_assignment

var AzureSecretStoreStackAssignmentCreateConfig = `
resource "snapcd_secret_store_stack_assignment" "this" { 
  secret_store_id = snapcd_azure_secret_store.this.id
  stack_id        = snapcd_stack.this.id
}`

var AzureSecretStoreNamespaceAssignmentCreateConfig = `
resource "snapcd_secret_store_namespace_assignment" "this" { 
  secret_store_id = snapcd_azure_secret_store.this.id
  namespace_id    = snapcd_namespace.this.id
}`

var AzureSecretStoreModuleAssignmentCreateConfig = `
resource "snapcd_secret_store_module_assignment" "this" { 
  secret_store_id = snapcd_azure_secret_store.this.id
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

var AwsSecretStoreStackAssignmentCreateConfig = `
resource "snapcd_secret_store_stack_assignment" "this" { 
  secret_store_id = snapcd_aws_secret_store.this.id
  stack_id        = snapcd_stack.this.id
}`

var AwsSecretStoreNamespaceAssignmentCreateConfig = `
resource "snapcd_secret_store_namespace_assignment" "this" { 
  secret_store_id = snapcd_aws_secret_store.this.id
  namespace_id    = snapcd_namespace.this.id
}`

var AwsSecretStoreModuleAssignmentCreateConfig = `
resource "snapcd_secret_store_module_assignment" "this" { 
  secret_store_id = snapcd_aws_secret_store.this.id
  module_id        = snapcd_module.this.id
}`
