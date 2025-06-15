// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/identity"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ResourceRoleAssignmentCreateConfig = `
resource "snapcd_resource_role_assignment" "this" { 
  resource_id  	 		  = snapcd_stack.this.id
  principal_id   		  = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  resource_discriminator  = "Stack"
  role_name 			  = "Owner"
}`

func TestAccResourceResourceRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + identity.ServicePrincipalDataSourceConfig + ResourceRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_resource_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceResourceRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + identity.ServicePrincipalDataSourceConfig + ResourceRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_resource_role_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_resource_role_assignment.this", "role_name", "Owner"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + identity.ServicePrincipalDataSourceConfig + `

resource "snapcd_resource_role_assignment" "this" { 
  resource_id  	 		  = snapcd_stack.this.id
  principal_id   		  = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  resource_discriminator  = "Stack"
  role_name 			  = "Contributor"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_resource_role_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_resource_role_assignment.this", "role_name", "Contributor"),
				),
			},
		},
	})
}

func TestAccResourceResourceRoleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + identity.ServicePrincipalDataSourceConfig + ResourceRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_resource_role_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_resource_role_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
