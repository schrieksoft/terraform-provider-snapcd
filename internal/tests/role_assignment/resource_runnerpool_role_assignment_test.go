// SPDX-License-Identifier: MPL-2.0

package role_assignment

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRunnerPoolRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.RunnerPoolCreateConfig + ServicePrincipalDataSourceConfig + RunnerPoolRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runnerpool_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerPoolRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.RunnerPoolCreateConfig + ServicePrincipalDataSourceConfig + RunnerPoolRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runnerpool_role_assignment.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.RunnerPoolCreateConfig + ServicePrincipalDataSourceConfig + RunnerPoolRoleAssignmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runnerpool_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerPoolRoleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.RunnerPoolCreateConfig + ServicePrincipalDataSourceConfig + RunnerPoolRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runnerpool_role_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_runnerpool_role_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
