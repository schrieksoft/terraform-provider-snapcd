// SPDX-License-Identifier: MPL-2.0

package role_assignment

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceOrganizationRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ServicePrincipalDataSourceConfig + OrganizationRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceOrganizationRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ServicePrincipalDataSourceConfig + OrganizationRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_role_assignment.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + ServicePrincipalDataSourceConfig + OrganizationRoleAssignmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceOrganizationRoleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ServicePrincipalDataSourceConfig + OrganizationRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_role_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_organization_role_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
