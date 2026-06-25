// SPDX-License-Identifier: MPL-2.0

package role_assignments

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceModuleRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + ServicePrincipalDataSourceConfig + ModuleRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + ServicePrincipalDataSourceConfig + ModuleRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_role_assignment.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + ServicePrincipalDataSourceConfig + ModuleRoleAssignmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleRoleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + ServicePrincipalDataSourceConfig + ModuleRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_role_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_role_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
