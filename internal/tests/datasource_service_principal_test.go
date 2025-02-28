// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceServicePrincipal(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + servicePrincipalCreateConfig + `

data "snapcd_service_principal" "this" {
	client_id = snapcd_service_principal.this.client_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_service_principal.this", "client_id", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
