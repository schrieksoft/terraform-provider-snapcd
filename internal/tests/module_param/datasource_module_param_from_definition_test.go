// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleParamFromDefinition(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromDefinitionCreateConfig + `
data "snapcd_module_param_from_definition" "this" {
	name 		= snapcd_module_param_from_definition.this.name
	module_id 	= snapcd_module_param_from_definition.this.module_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_param_from_definition.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_param_from_definition.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
