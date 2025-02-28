// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceRunnerPool(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + runnerPoolCreateConfig + `

data "snapcd_runner_pool" "this" {
	name = snapcd_runner_pool.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_runner_pool.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_runner_pool.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
