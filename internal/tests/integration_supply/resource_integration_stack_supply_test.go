// SPDX-License-Identifier: MPL-2.0

package integration_supply

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceIntegrationStackSupply_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + IntegrationDataSourceConfig + IntegrationStackSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_integration_stack_supply.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceIntegrationStackSupply_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + IntegrationDataSourceConfig + IntegrationStackSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_integration_stack_supply.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_integration_stack_supply.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
