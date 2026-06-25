// SPDX-License-Identifier: MPL-2.0

package integration

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceIntegration(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + IntegrationDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_integration.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_integration.this", "name", "debug-slack"),
					resource.TestCheckResourceAttr("data.snapcd_integration.this", "integration_type", "Slack"),
				),
			},
		},
	})
}
