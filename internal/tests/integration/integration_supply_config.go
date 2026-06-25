package integration

var IntegrationStackSupplyCreateConfig = `
resource "snapcd_integration_stack_supply" "this" {
  integration_id = data.snapcd_integration.this.id
  stack_id       = snapcd_stack.this.id
}`

var IntegrationNamespaceSupplyCreateConfig = `
resource "snapcd_integration_namespace_supply" "this" {
  integration_id = data.snapcd_integration.this.id
  namespace_id   = snapcd_namespace.this.id
}`

var IntegrationModuleSupplyCreateConfig = `
resource "snapcd_integration_module_supply" "this" {
  integration_id = data.snapcd_integration.this.id
  module_id      = snapcd_module.this.id
}`
