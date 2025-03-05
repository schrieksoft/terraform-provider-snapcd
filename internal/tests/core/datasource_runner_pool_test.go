// SPDX-License-Identifier: MPL-2.0

package core

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceRunnerPool(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + RunnerPoolCreateConfig + `

data "snapcd_runner_pool" "this" {
	name = snapcd_runner_pool.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_runner_pool.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_runner_pool.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
