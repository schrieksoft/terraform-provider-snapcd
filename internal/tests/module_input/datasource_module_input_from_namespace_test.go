// SPDX-License-Identifier: MPL-2.0

package module_input

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleInputFromNamespace(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + NamespaceInputFromLiteralCreateConfig + ModuleInputFromNamespaceCreateConfig + `
data "snapcd_module_input_from_namespace" "this" {
	name 		= snapcd_module_input_from_namespace.this.name
	module_id 	= snapcd_module_input_from_namespace.this.module_id
	input_kind  = "Param"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_input_from_namespace.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_input_from_namespace.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
