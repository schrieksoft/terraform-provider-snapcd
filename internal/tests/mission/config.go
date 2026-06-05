// SPDX-License-Identifier: MPL-2.0

package mission

var OrganizationMissionCreateConfig = `
resource "snapcd_organization_mission" "this" {
  agent_id     = snapcd_agent.this.id
  mission_type = "AutoDiagnose"
}`

var StackMissionCreateConfig = `
resource "snapcd_stack_mission" "this" {
  agent_id     = snapcd_agent.this.id
  stack_id     = snapcd_stack.this.id
  mission_type = "AutoDiagnose"
}`

var NamespaceMissionCreateConfig = `
resource "snapcd_namespace_mission" "this" {
  agent_id     = snapcd_agent.this.id
  namespace_id = snapcd_namespace.this.id
  mission_type = "AutoDiagnose"
}`

var ModuleMissionCreateConfig = `
resource "snapcd_module_mission" "this" {
  agent_id     = snapcd_agent.this.id
  module_id    = snapcd_module.this.id
  mission_type = "AutoDiagnose"
}`
