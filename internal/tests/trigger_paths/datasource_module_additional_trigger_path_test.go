// SPDX-License-Identifier: MPL-2.0

package trigger_paths

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleAdditionalTriggerPath(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleAdditionalTriggerPathCreateConfig + `
data "snapcd_module_additional_trigger_path" "this" {
	path      = snapcd_module_additional_trigger_path.this.path
	module_id = snapcd_module_additional_trigger_path.this.module_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_additional_trigger_path.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_additional_trigger_path.this", "path", providerconfig.AppendRandomString("shared/scripts%s")),
				),
			},
		},
	})
}
