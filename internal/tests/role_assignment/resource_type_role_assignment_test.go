// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/identity"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var TypeRoleAssignmentCreateConfig = `
resource "snapcd_type_role_assignment" "this" { 
  principal_id   		  = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  resource_discriminator  = "Stack"
  role_name 			  = "Owner"
}`

func TestAccResourceTypeRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + identity.ServicePrincipalDataSourceConfig + TypeRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_type_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceTypeRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + identity.ServicePrincipalDataSourceConfig + TypeRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_type_role_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_type_role_assignment.this", "role_name", "Owner"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + identity.ServicePrincipalDataSourceConfig + `

resource "snapcd_type_role_assignment" "this" { 
  principal_id   		  = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  resource_discriminator  = "Stack"
  role_name 			  = "Contributor"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_type_role_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_type_role_assignment.this", "role_name", "Contributor"),
				),
			},
		},
	})
}

func TestAccResourceTypeRoleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + identity.ServicePrincipalDataSourceConfig + TypeRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_type_role_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_type_role_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
