// SPDX-License-Identifier: MPL-2.0

package policies

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModulePulumiInlinePolicy(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModulePulumiInlinePolicyCreateConfig + `
data "snapcd_module_pulumi_inline_policy" "this" {
	name      = snapcd_module_pulumi_inline_policy.this.name
	module_id = snapcd_module_pulumi_inline_policy.this.module_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_pulumi_inline_policy.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_pulumi_inline_policy.this", "name", providerconfig.AppendRandomString("mypolicy%s")),
					resource.TestCheckResourceAttr("data.snapcd_module_pulumi_inline_policy.this", "evaluate_on", "ApplyAndDestroy"),
				),
			},
		},
	})
}
