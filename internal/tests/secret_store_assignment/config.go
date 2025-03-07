package secret_store

var SecretStoreStackAssignmentCreateConfig = `
resource "snapcd_secret_store_stack_assignment" "this" { 
  permission  		= "ReadWrite"
  secret_store_id = snapcd_azure_key_vault_secret_store.this.id
  stack_id        = snapcd_stack.this.id
}`

var SecretStoreNamespaceAssignmentCreateConfig = `
resource "snapcd_secret_store_namespace_assignment" "this" { 
  permission  		= "ReadWrite"
  secret_store_id = snapcd_azure_key_vault_secret_store.this.id
  namespace_id    = snapcd_namespace.this.id
}`

var SecretStoreModuleAssignmentCreateConfig = `
resource "snapcd_secret_store_module_assignment" "this" { 
  permission  		= "ReadWrite"
  secret_store_id = snapcd_azure_key_vault_secret_store.this.id
  module_id        = snapcd_module.this.id
}`
