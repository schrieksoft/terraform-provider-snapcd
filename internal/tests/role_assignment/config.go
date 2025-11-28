package role_assignment

var ServicePrincipalDataSourceConfig = `
data "snapcd_service_principal" "this" {
	client_id = "debugTestTarget1"
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
  stack_id                = snapcd_stack.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Owner"
}
`

var StackRoleAssignmentUpdateConfig = `
resource "snapcd_stack_role_assignment" "this" {
  stack_id                = snapcd_stack.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Contributor"
}
`

// Namespace Role Assignment Configs
var NamespaceRoleAssignmentCreateConfig = `
resource "snapcd_namespace_role_assignment" "this" {
  namespace_id            = snapcd_namespace.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Owner"
}
`

var NamespaceRoleAssignmentUpdateConfig = `
resource "snapcd_namespace_role_assignment" "this" {
  namespace_id            = snapcd_namespace.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Reader"
}
`

// Module Role Assignment Configs
var ModuleRoleAssignmentCreateConfig = `
resource "snapcd_module_role_assignment" "this" {
  module_id               = snapcd_module.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Owner"
}
`

var ModuleRoleAssignmentUpdateConfig = `
resource "snapcd_module_role_assignment" "this" {
  module_id               = snapcd_module.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "IdentityAccessManager"
}
`

// Runner Role Assignment Configs
var RunnerRoleAssignmentCreateConfig = `
resource "snapcd_runner_role_assignment" "this" {
  runner_id          	 = snapcd_runner.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Owner"
}
`

var RunnerRoleAssignmentUpdateConfig = `
resource "snapcd_runner_role_assignment" "this" {
  runner_id          	  = snapcd_runner.this.id
  principal_id            = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name               = "Reader"
}
`
