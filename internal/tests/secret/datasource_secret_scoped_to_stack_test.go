// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSecretScopedToStack(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + `

data "snapcd_secret_scoped_to_stack" "this" {
	name 	  = "debug"
    stack_id  = "10000000-0000-0000-0000-000000000000"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_secret_scoped_to_stack.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_secret_scoped_to_stack.this", "id", "10000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}
