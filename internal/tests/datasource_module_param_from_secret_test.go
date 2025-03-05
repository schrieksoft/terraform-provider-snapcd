// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleParamFromSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromSecretCreateConfig + `
data "snapcd_module_param_from_secret" "this" {
	name 		= snapcd_module_param_from_secret.this.name
	module_id 	= snapcd_module_param_from_secret.this.module_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_param_from_secret.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_param_from_secret.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
