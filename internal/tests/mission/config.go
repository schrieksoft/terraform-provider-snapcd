// SPDX-License-Identifier: MPL-2.0

package mission

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
)

var OrganizationMissionCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_organization_mission" "this" {
  agent_id     = snapcd_agent.this.id
  name         = "somevalue%s"
  mission_type = "AutoDiagnose"
}`)

var StackMissionCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_stack_mission" "this" {
  agent_id     = snapcd_agent.this.id
  stack_id     = snapcd_stack.this.id
  name         = "somevalue%s"
  mission_type = "AutoDiagnose"
}`)

var NamespaceMissionCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_mission" "this" {
  agent_id     = snapcd_agent.this.id
  namespace_id = snapcd_namespace.this.id
  name         = "somevalue%s"
  mission_type = "AutoDiagnose"
}`)

var ModuleMissionCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_mission" "this" {
  agent_id     = snapcd_agent.this.id
  module_id    = snapcd_module.this.id
  name         = "somevalue%s"
  mission_type = "AutoDiagnose"
}`)
