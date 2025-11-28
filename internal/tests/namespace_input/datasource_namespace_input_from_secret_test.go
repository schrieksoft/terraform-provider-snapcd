// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceInputFromSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfig + `
data "snapcd_namespace_input_from_secret" "this" {
	name 		 = snapcd_namespace_input_from_secret.this.name
	namespace_id = snapcd_namespace_input_from_secret.this.namespace_id
	input_kind 	 = "Param"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_input_from_secret.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_input_from_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
