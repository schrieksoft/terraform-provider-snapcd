// SPDX-License-Identifier: MPL-2.0

package missions

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceOrganizationMission_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.AgentCreateConfig + OrganizationMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_mission.this", "id"),
					resource.TestCheckResourceAttr("snapcd_organization_mission.this", "mission_type", "AutoDiagnose"),
				),
			},
		},
	})
}

func TestAccResourceOrganizationMission_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.AgentCreateConfig + OrganizationMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_mission.this", "id"),
					resource.TestCheckResourceAttr("snapcd_organization_mission.this", "mission_type", "AutoDiagnose"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + testdata.AgentCreateConfig + `
resource "snapcd_organization_mission" "this" {
  agent_id     = snapcd_agent.this.id
  mission_type = "ApprovalRecommend"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_mission.this", "id"),
					resource.TestCheckResourceAttr("snapcd_organization_mission.this", "mission_type", "ApprovalRecommend"),
				),
			},
		},
	})
}

func TestAccResourceOrganizationMission_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.AgentCreateConfig + OrganizationMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_mission.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_organization_mission.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
