// SPDX-License-Identifier: MPL-2.0

package role_assignments

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStackRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.StackCreateConfig + ServicePrincipalDataSourceConfig + StackRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceStackRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.StackCreateConfig + ServicePrincipalDataSourceConfig + StackRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_role_assignment.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.StackCreateConfig + ServicePrincipalDataSourceConfig + StackRoleAssignmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceStackRoleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.StackCreateConfig + ServicePrincipalDataSourceConfig + StackRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_role_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_stack_role_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
