// SPDX-License-Identifier: MPL-2.0

package state_store

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceStateStore(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + StateStoreCreateConfig + `

data "snapcd_state_store" "this" {
	name = snapcd_state_store.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_state_store.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_state_store.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
