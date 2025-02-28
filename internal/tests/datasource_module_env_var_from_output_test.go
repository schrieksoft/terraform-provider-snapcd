// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleEnvVarFromOutput(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromOutputCreateConfig + `
data "snapcd_module_env_var_from_output" "this" {
	name 		= snapcd_module_env_var_from_output.this.name
	module_id 	= snapcd_module_env_var_from_output.this.module_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_env_var_from_output.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_env_var_from_output.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
