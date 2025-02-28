// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleParamFromNamespace(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromNamespaceCreateConfig + `
data "snapcd_module_param_from_namespace" "this" {
	name 		= snapcd_module_param_from_namespace.this.name
	module_id 	= snapcd_module_param_from_namespace.this.module_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_param_from_namespace.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_param_from_namespace.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
