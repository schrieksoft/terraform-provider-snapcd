// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)


func TestAccResourceGroupMember_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + GroupCreateConfig + ServicePrincipalCreateConfig + GroupMemberCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group_member.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceGroupMember_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + GroupCreateConfig + ServicePrincipalCreateConfig + GroupMemberCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group_member.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + GroupCreateConfig + ServicePrincipalCreateConfig + providerconfig.AppendRandomString(`

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
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + GroupCreateConfig + ServicePrincipalCreateConfig + GroupMemberCreateConfig,
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
