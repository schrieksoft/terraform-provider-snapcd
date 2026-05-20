// SPDX-License-Identifier: MPL-2.0

package mission

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceModuleMission_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + core.AgentCreateConfig + ModuleMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_mission.this", "id"),
					resource.TestCheckResourceAttrSet("snapcd_module_mission.this", "module_id"),
					resource.TestCheckResourceAttr("snapcd_module_mission.this", "mission_type", "AutoDiagnose"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + core.AgentCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_mission" "this" {
  agent_id     = snapcd_agent.this.id
  module_id    = snapcd_module.this.id
  name         = "somevalue%s"
  mission_type = "SplitMonolithicState"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snapcd_module_mission.this", "mission_type", "SplitMonolithicState"),
				),
			},
		},
	})
}

func TestAccResourceModuleMission_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + core.AgentCreateConfig + ModuleMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_mission.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_mission.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
