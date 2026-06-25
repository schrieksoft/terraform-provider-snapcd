package runner

var RunnerStackSupplyCreateConfig = `
resource "snapcd_runner_stack_supply" "this" { 
  runner_id = snapcd_runner.this.id
  stack_id        = snapcd_stack.this.id
}`

var RunnerNamespaceSupplyCreateConfig = `
resource "snapcd_runner_namespace_supply" "this" { 
  runner_id = snapcd_runner.this.id
  namespace_id    = snapcd_namespace.this.id
}`

var RunnerModuleSupplyCreateConfig = `
resource "snapcd_runner_module_supply" "this" { 
  runner_id = snapcd_runner.this.id
  module_id        = snapcd_module.this.id
}`

var RunnerServicePrincipalSupplyCreateConfig = `
data "snapcd_service_principal" "this" {
	client_id = "debugTestTarget1"
}

resource "snapcd_runner_service_principal_supply" "this" { 
  runner_id        = snapcd_runner.this.id
  service_principal_id  = data.snapcd_service_principal.this.id
}`
