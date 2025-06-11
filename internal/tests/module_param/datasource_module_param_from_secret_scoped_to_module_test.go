// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleParamFromSecretScopedToModule(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleParamFromSecretScopedToModuleCreateConfig + `
data "snapcd_module_param_from_secret_scoped_to_module" "this" {
	secret_scoped_to_module_id 		= snapcd_module_param_from_secret_scoped_to_module.this.secret_scoped_to_module_id
	module_id 	                    = snapcd_module_param_from_secret_scoped_to_module.this.module_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_param_from_secret_scoped_to_module.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_param_from_secret_scoped_to_module.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
