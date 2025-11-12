package runner_assignment

var RunnerStackAssignmentCreateConfig = `
resource "snapcd_runner_stack_assignment" "this" { 
  runner_id = snapcd_runner.this.id
  stack_id        = snapcd_stack.this.id
}`

var RunnerNamespaceAssignmentCreateConfig = `
resource "snapcd_runner_namespace_assignment" "this" { 
  runner_id = snapcd_runner.this.id
  namespace_id    = snapcd_namespace.this.id
}`

var RunnerModuleAssignmentCreateConfig = `
resource "snapcd_runner_module_assignment" "this" { 
  runner_id = snapcd_runner.this.id
  module_id        = snapcd_module.this.id
}`

var RunnerServicePrincipalAssignmentCreateConfig = `
data "snapcd_service_principal" "this" {
	client_id = "debugTestTarget1"
}

resource "snapcd_runner_service_principal_assignment" "this" { 
  runner_id        = snapcd_runner.this.id
  service_principal_id  = data.snapcd_service_principal.this.id
}`
