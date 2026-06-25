package agent

var AgentStackSupplyCreateConfig = `
resource "snapcd_agent_stack_supply" "this" {
  agent_id = snapcd_agent.this.id
  stack_id = snapcd_stack.this.id
}`

var AgentNamespaceSupplyCreateConfig = `
resource "snapcd_agent_namespace_supply" "this" {
  agent_id     = snapcd_agent.this.id
  namespace_id = snapcd_namespace.this.id
}`

var AgentModuleSupplyCreateConfig = `
resource "snapcd_agent_module_supply" "this" {
  agent_id  = snapcd_agent.this.id
  module_id = snapcd_module.this.id
}`
