// SPDX-License-Identifier: MPL-2.0

package missions

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceNamespaceMission_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + testdata.AgentCreateConfig + NamespaceMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_mission.this", "id"),
					resource.TestCheckResourceAttrSet("snapcd_namespace_mission.this", "namespace_id"),
					resource.TestCheckResourceAttr("snapcd_namespace_mission.this", "mission_type", "AutoDiagnose"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceMission_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + testdata.AgentCreateConfig + NamespaceMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_mission.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_mission.this", "mission_type", "AutoDiagnose"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + testdata.AgentCreateConfig + `
resource "snapcd_namespace_mission" "this" {
  agent_id     = snapcd_agent.this.id
  namespace_id = snapcd_namespace.this.id
  mission_type = "SummarizeJob"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snapcd_namespace_mission.this", "mission_type", "SummarizeJob"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceMission_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + testdata.AgentCreateConfig + NamespaceMissionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_mission.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_mission.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
