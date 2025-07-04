// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleInputFromOutput(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + core.ModuleCreateConfigDeltaTwo + ModuleInputFromOutputCreateConfig + `
data "snapcd_module_input_from_output" "this" {
	name 		= snapcd_module_input_from_output.this.name
	module_id 	= snapcd_module_input_from_output.this.module_id
	input_kind = "Param"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_input_from_output.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_input_from_output.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
