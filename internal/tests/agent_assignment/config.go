package agent_assignment

var AgentStackAssignmentCreateConfig = `
resource "snapcd_agent_stack_assignment" "this" {
  agent_id = snapcd_agent.this.id
  stack_id = snapcd_stack.this.id
}`

var AgentNamespaceAssignmentCreateConfig = `
resource "snapcd_agent_namespace_assignment" "this" {
  agent_id     = snapcd_agent.this.id
  namespace_id = snapcd_namespace.this.id
}`

var AgentModuleAssignmentCreateConfig = `
resource "snapcd_agent_module_assignment" "this" {
  agent_id  = snapcd_agent.this.id
  module_id = snapcd_module.this.id
}`
