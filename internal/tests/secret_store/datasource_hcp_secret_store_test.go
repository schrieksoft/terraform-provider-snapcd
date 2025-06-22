// SPDX-License-Identifier: MPL-2.0

package secret_store

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceHcpSecretStore(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + HcpSecretStoreCreateConfig + `

data "snapcd_hcp_secret_store" "this" {
	name = snapcd_hcp_secret_store.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_hcp_secret_store.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_hcp_secret_store.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
