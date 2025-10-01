// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceInputFromLiteral(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromLiteralCreateConfig + `
data "snapcd_namespace_input_from_literal" "this" {
	name 		 = snapcd_namespace_input_from_literal.this.name
	namespace_id = snapcd_namespace_input_from_literal.this.namespace_id
	input_kind 	 = "Param"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_input_from_literal.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_input_from_literal.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
