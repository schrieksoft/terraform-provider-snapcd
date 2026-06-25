// SPDX-License-Identifier: MPL-2.0

package namespace_input

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceInputFromDefinition(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceInputFromDefinitionCreateConfig + `
data "snapcd_namespace_input_from_definition" "this" {
	name 		 = snapcd_namespace_input_from_definition.this.name
	namespace_id = snapcd_namespace_input_from_definition.this.namespace_id
	input_kind 	 = "Param"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_input_from_definition.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_input_from_definition.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
