// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + `

data "snapcd_namespace_secret" "this" {
	name 		 = "debug"
    namespace_id = "99999999-9999-9999-9999-999999999999"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_secret.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_secret.this", "id", "99999999-9999-9999-9999-999999999902"),
				),
			},
		},
	})
}
