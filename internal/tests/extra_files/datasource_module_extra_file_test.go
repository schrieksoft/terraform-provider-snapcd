// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleExtraFile(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleExtraFileCreateConfig + `
data "snapcd_module_extra_file" "this" {
	file_name = snapcd_module_extra_file.this.file_name
	module_id = snapcd_module_extra_file.this.module_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_extra_file.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_extra_file.this", "file_name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
