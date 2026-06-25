// SPDX-License-Identifier: MPL-2.0

package role_assignments

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRunnerRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.RunnerCreateConfig + ServicePrincipalDataSourceConfig + RunnerRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_role_assignment.this", "id"),
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
				Config: providerconfig.ProviderConfig() + testdata.RunnerCreateConfig + ServicePrincipalDataSourceConfig + RunnerRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_role_assignment.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.RunnerCreateConfig + ServicePrincipalDataSourceConfig + RunnerRoleAssignmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_role_assignment.this", "id"),
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
				Config: providerconfig.ProviderConfig() + testdata.RunnerCreateConfig + ServicePrincipalDataSourceConfig + RunnerRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_role_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_runner_role_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
