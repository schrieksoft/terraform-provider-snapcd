// SPDX-License-Identifier: MPL-2.0

package mission

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStackMission_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + core.AgentCreateConfig + StackMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_mission.this", "id"),
					resource.TestCheckResourceAttrSet("snapcd_stack_mission.this", "stack_id"),
					resource.TestCheckResourceAttr("snapcd_stack_mission.this", "mission_type", "AutoDiagnose"),
				),
			},
		},
	})
}

func TestAccResourceStackMission_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + core.AgentCreateConfig + StackMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_mission.this", "id"),
					resource.TestCheckResourceAttr("snapcd_stack_mission.this", "mission_type", "AutoDiagnose"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + core.AgentCreateConfig + `
resource "snapcd_stack_mission" "this" {
  agent_id     = snapcd_agent.this.id
  stack_id     = snapcd_stack.this.id
  mission_type = "ApprovalRecommend"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snapcd_stack_mission.this", "mission_type", "ApprovalRecommend"),
				),
			},
		},
	})
}

func TestAccResourceStackMission_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + core.AgentCreateConfig + StackMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_mission.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_stack_mission.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
