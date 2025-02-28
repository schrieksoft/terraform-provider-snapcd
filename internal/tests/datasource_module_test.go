// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModule(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + `
data "snapcd_module" "this" {
	name 			  = snapcd_module.this.name
	namespace_id 	  = snapcd_namespace.this.id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
