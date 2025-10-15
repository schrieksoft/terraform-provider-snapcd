package role_assignment

var ServicePrincipalDataSourceConfig = `
data "snapcd_service_principal" "this" {
	client_id = "debug"
}
`

var UserDataSourceConfig = `
data "snapcd_user" "this" {
  user_name  = "kschriek@gmail.com"
}
`

var StackDataSourceConfig = `
data "snapcd_stack" "this" {
  name = "debug"
}
`

var NamespaceDataSourceConfig = `
data "snapcd_namespace" "this" {
  stack_id = data.snapcd_stack.this.id
  name     = "debug"
}
`

var ModuleDataSourceConfig = `
data "snapcd_module" "this" {
  namespace_id = data.snapcd_namespace.this.id
  name         = "debug"
}
`

// Organization Role Assignment Configs
var OrganizationRoleAssignmentCreateConfig = `
resource "snapcd_organization_role_assignment" "this" {
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Reader"
}
`

var OrganizationRoleAssignmentUpdateConfig = `
resource "snapcd_organization_role_assignment" "this" {
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "IdentityAccessManager"
}
`

// Stack Role Assignment Configs
var StackRoleAssignmentCreateConfig = `
resource "snapcd_stack_role_assignment" "this" {
  stack_id                = data.snapcd_stack.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Owner"
}
`

var StackRoleAssignmentUpdateConfig = `
resource "snapcd_stack_role_assignment" "this" {
  stack_id                = data.snapcd_stack.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Contributor"
}
`

// Namespace Role Assignment Configs
var NamespaceRoleAssignmentCreateConfig = `
resource "snapcd_namespace_role_assignment" "this" {
  namespace_id            = data.snapcd_namespace.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Owner"
}
`

var NamespaceRoleAssignmentUpdateConfig = `
resource "snapcd_namespace_role_assignment" "this" {
  namespace_id            = data.snapcd_namespace.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Reader"
}
`

// Module Role Assignment Configs
var ModuleRoleAssignmentCreateConfig = `
resource "snapcd_module_role_assignment" "this" {
  module_id               = data.snapcd_module.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Owner"
}
`

var ModuleRoleAssignmentUpdateConfig = `
resource "snapcd_module_role_assignment" "this" {
  module_id               = data.snapcd_module.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "IdentityAccessManager"
}
`
