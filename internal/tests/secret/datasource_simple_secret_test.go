// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSimpleSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + SimpleSecretCreateConfig + `

data "snapcd_simple_secret" "this" {
	name = snapcd_simple_secret.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_simple_secret.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_simple_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
