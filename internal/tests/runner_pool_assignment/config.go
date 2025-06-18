package runner_pool_assignment

var RunnerPoolStackAssignmentCreateConfig = `
resource "snapcd_runner_pool_stack_assignment" "this" { 
  runner_pool_id = snapcd_runner_pool.this.id
  stack_id        = snapcd_stack.this.id
}`

var RunnerPoolNamespaceAssignmentCreateConfig = `
resource "snapcd_runner_pool_namespace_assignment" "this" { 
  runner_pool_id = snapcd_runner_pool.this.id
  namespace_id    = snapcd_namespace.this.id
}`

var RunnerPoolModuleAssignmentCreateConfig = `
resource "snapcd_runner_pool_module_assignment" "this" { 
  runner_pool_id = snapcd_runner_pool.this.id
  module_id        = snapcd_module.this.id
}`

var RunnerPoolServicePrincipalAssignmentCreateConfig = `
data "snapcd_service_principal" "this" {
	client_id = "IntegratedRunner"
}

resource "snapcd_runner_pool_service_principal_assignment" "this" { 
  runner_pool_id        = snapcd_runner_pool.this.id
  service_principal_id  = data.snapcd_service_principal.this.id
}`
