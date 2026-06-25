// SPDX-License-Identifier: MPL-2.0

package integration

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceIntegrationModuleSupply_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + IntegrationDataSourceConfig + IntegrationModuleSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_integration_module_supply.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceIntegrationModuleSupply_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + IntegrationDataSourceConfig + IntegrationModuleSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_integration_module_supply.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_integration_module_supply.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
