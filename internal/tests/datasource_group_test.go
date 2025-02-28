// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGroup(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + groupCreateConfig + `

data "snapcd_group" "this" {
	name = snapcd_group.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_group.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_group.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
