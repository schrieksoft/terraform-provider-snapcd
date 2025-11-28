// SPDX-License-Identifier: MPL-2.0

package core

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceRunner(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + RunnerCreateConfig + `

data "snapcd_runner" "this" {
	name = snapcd_runner.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_runner.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_runner.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
