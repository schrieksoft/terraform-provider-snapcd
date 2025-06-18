// SPDX-License-Identifier: MPL-2.0

package core

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRunnerPool_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + RunnerPoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerPool_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + RunnerPoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool.this", "id"),
					resource.TestCheckResourceAttr("snapcd_runner_pool.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_runner_pool" "this" { 
  name = "someNEWvalue%s"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool.this", "id"),
					resource.TestCheckResourceAttr("snapcd_runner_pool.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceRunnerPool_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + RunnerPoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_runner_pool.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
