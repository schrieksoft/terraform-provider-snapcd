package secret

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret_store"
	"terraform-provider-snapcd/internal/tests/secret_store_assignment"
)

var SecretScopedToNamespaceCreateConfig = core.NamespaceCreateConfig + secret_store.AwsSecretStoreCreateConfig + secret_store_assignment.AwsSecretStoreNamespaceAssignmentCreateConfig + SecretScopedToNamespaceCreateConfigDelta

var SecretScopedToNamespaceCreateConfigDelta = providerconfig.AppendRandomString(`
resource "snapcd_secret_scoped_to_namespace" "this" { 
  depends_on         = [snapcd_secret_store_namespace_assignment.this]
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_aws_secret_store.this.id
  namespace_id 	     = snapcd_namespace.this.id
}`)

var SecretScopedToModuleCreateConfig = core.ModuleCreateConfig + secret_store.AwsSecretStoreCreateConfig + secret_store_assignment.AwsSecretStoreModuleAssignmentCreateConfig + SecretScopedToModuleCreateConfigDelta

var SecretScopedToModuleCreateConfigDelta = providerconfig.AppendRandomString(`
resource "snapcd_secret_scoped_to_module" "this" { 
  depends_on         = [snapcd_secret_store_module_assignment.this]
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_aws_secret_store.this.id
  module_id 	       = snapcd_module.this.id
}`)

var SecretScopedToStackCreateConfig = core.StackCreateConfig + secret_store.AwsSecretStoreCreateConfig + secret_store_assignment.AwsSecretStoreStackAssignmentCreateConfig + SecretScopedToStackCreateConfigDelta

var SecretScopedToStackCreateConfigDelta = providerconfig.AppendRandomString(`
resource "snapcd_secret_scoped_to_stack" "this" { 
  depends_on         = [snapcd_secret_store_stack_assignment.this]
  name  		         = "somevalue%s"
  remote_secret_name = "name-in-remote"
  secret_store_id    = snapcd_aws_secret_store.this.id
  stack_id 	         = snapcd_stack.this.id
}`)
