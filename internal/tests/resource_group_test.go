// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var groupCreateConfig = appendRandomString(`
resource "snapcd_group" "this" { 
  name  = "somevalue%s"
}`)

func TestAccResourceGroup_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + groupCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceGroup_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + groupCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group.this", "id"),
					resource.TestCheckResourceAttr("snapcd_group.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + appendRandomString(`
resource "snapcd_group" "this" { 
  name = "someNEWvalue%s"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group.this", "id"),
					resource.TestCheckResourceAttr("snapcd_group.this", "name", appendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceGroup_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + groupCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_group.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
