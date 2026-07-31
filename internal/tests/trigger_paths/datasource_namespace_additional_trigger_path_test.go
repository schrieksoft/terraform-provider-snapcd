// SPDX-License-Identifier: MPL-2.0

package trigger_paths

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceAdditionalTriggerPath(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceAdditionalTriggerPathCreateConfig + `
data "snapcd_namespace_additional_trigger_path" "this" {
	path         = snapcd_namespace_additional_trigger_path.this.path
	namespace_id = snapcd_namespace_additional_trigger_path.this.namespace_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_additional_trigger_path.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_additional_trigger_path.this", "path", providerconfig.AppendRandomString("shared/config%s")),
				),
			},
		},
	})
}
