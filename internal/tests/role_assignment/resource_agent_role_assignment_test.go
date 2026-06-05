// SPDX-License-Identifier: MPL-2.0

package role_assignment

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAgentRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.AgentCreateConfig + ServicePrincipalDataSourceConfig + AgentRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceAgentRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.AgentCreateConfig + ServicePrincipalDataSourceConfig + AgentRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_role_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_agent_role_assignment.this", "role_name", "Owner"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.AgentCreateConfig + ServicePrincipalDataSourceConfig + AgentRoleAssignmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_role_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_agent_role_assignment.this", "role_name", "Reader"),
				),
			},
		},
	})
}

func TestAccResourceAgentRoleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.AgentCreateConfig + ServicePrincipalDataSourceConfig + AgentRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_role_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_agent_role_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
