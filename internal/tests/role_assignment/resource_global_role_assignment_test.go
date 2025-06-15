// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/identity"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const GlobalRoleAssignmentCreateConfig = `
resource "snapcd_global_role_assignment" "this" {
  principal_id   		  = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name 			  = "Administrator"
}`

func TestAccResourceGlobalRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + identity.ServicePrincipalDataSourceConfig + GlobalRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_global_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceGlobalRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + identity.ServicePrincipalDataSourceConfig + GlobalRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_global_role_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_global_role_assignment.this", "role_name", "Administrator"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + identity.ServicePrincipalDataSourceConfig + `

resource "snapcd_global_role_assignment" "this" {
  principal_id   		  = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  role_name 			  = "IdentityAccessAdministrator"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_global_role_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_global_role_assignment.this", "role_name", "IdentityAccessAdministrator"),
				),
			},
		},
	})
}

func TestAccResourceGlobalRoleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + identity.ServicePrincipalDataSourceConfig + GlobalRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_global_role_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_global_role_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
