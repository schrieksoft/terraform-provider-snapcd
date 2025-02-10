// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var resourceRoleAssignmentCreateConfig = `
resource "snapcd_resource_role_assignment" "this" { 
  resource_id  	 		  = snapcd_stack.this.id
  principal_id   		  = snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
  resource_discriminator  = "Stack"
  role_name 			  = "Owner"
}`

func TestAccResourceResourceRoleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + stackCreateConfig + servicePrincipalCreateConfig + resourceRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_resource_role_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceResourceRoleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + stackCreateConfig + servicePrincipalCreateConfig + resourceRoleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_resource_role_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_resource_role_assignment.this", "role_name", "Owner"),
				),
			},
			{
				Config: providerConfig + stackCreateConfig + servicePrincipalCreateConfig + `

resource "snapcd_resource_role_assignment" "this" { 
  resource_id  	 		  = snapcd_stack.this.id
  principal_id   		  = snapcd_service_principal.this.id
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + stackCreateConfig + servicePrincipalCreateConfig + resourceRoleAssignmentCreateConfig,
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
