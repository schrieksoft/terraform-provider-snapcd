// SPDX-License-Identifier: MPL-2.0

package role_assignment

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRunnerRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.RunnerCreateConfig + ServicePrincipalDataSourceConfig + RunnerRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runnerpool_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.RunnerCreateConfig + ServicePrincipalDataSourceConfig + RunnerRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runnerpool_role_assignment.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.RunnerCreateConfig + ServicePrincipalDataSourceConfig + RunnerRoleAssignmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runnerpool_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerRoleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.RunnerCreateConfig + ServicePrincipalDataSourceConfig + RunnerRoleAssignmentCreateConfig,
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
