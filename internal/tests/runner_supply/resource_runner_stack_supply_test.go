// SPDX-License-Identifier: MPL-2.0

package runner_supply

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRunnerStackSupply_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + core.RunnerCreateConfig + RunnerStackSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_stack_supply.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerStackSupply_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + core.RunnerCreateConfig + RunnerStackSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_stack_supply.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_runner_stack_supply.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
