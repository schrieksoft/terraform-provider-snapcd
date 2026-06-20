// SPDX-License-Identifier: MPL-2.0

package runner_supply

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRunnerNamespaceSupply_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + core.RunnerCreateConfig + RunnerNamespaceSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_namespace_supply.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerNamespaceSupply_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + core.RunnerCreateConfig + RunnerNamespaceSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_namespace_supply.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_runner_namespace_supply.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
