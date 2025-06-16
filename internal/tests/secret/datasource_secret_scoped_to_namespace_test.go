// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSecretScopedToNamespace(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + SecretScopedToNamespaceCreateConfig + `

data "snapcd_secret_scoped_to_namespace" "this" {
	name = snapcd_secret_scoped_to_namespace.this.name
    namespace_id = snapcd_namespace.this.id

}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_secret_scoped_to_namespace.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_secret_scoped_to_namespace.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
