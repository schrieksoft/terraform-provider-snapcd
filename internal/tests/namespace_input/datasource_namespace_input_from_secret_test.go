// SPDX-License-Identifier: MPL-2.0

package namespace_input

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceInputFromSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfig + `
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
