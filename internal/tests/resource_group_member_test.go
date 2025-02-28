// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const groupMemberCreateConfig = `
resource "snapcd_group_member" "this" { 
  group_id  	 		  = snapcd_group.this.id
  principal_id   		  = snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
}`

func TestAccResourceGroupMember_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + groupCreateConfig + servicePrincipalCreateConfig + groupMemberCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group_member.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceGroupMember_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + groupCreateConfig + servicePrincipalCreateConfig + groupMemberCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group_member.this", "id"),
				),
			},
			{
				Config: providerConfig + groupCreateConfig + servicePrincipalCreateConfig + appendRandomString(`

resource "snapcd_service_principal" "new" { 
  client_id  	 = "someNEWvalue%s"
  client_secret  = "veryverysecret"
  scopes    	 = ["foo","bar","ban", "baz"]
  display_name   = "foo"
}

resource "snapcd_group_member" "this" { 
  group_id  	 		  = snapcd_group.this.id
  principal_id   		  = snapcd_service_principal.new.id
  principal_discriminator = "ServicePrincipal"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group_member.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceGroupMember_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + groupCreateConfig + servicePrincipalCreateConfig + groupMemberCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group_member.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_group_member.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
