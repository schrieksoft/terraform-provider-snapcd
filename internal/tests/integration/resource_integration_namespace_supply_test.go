// SPDX-License-Identifier: MPL-2.0

package integration

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceIntegrationNamespaceSupply_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + IntegrationDataSourceConfig + IntegrationNamespaceSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_integration_namespace_supply.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceIntegrationNamespaceSupply_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + IntegrationDataSourceConfig + IntegrationNamespaceSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_integration_namespace_supply.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_integration_namespace_supply.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
