// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceServicePrincipal(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + `

data "snapcd_service_principal" "this" {
	client_id = "DefaultRunner"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_service_principal.this", "client_id", "DefaultRunner"),
				),
			},
		},
	})
}
