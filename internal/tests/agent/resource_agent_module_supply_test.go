// SPDX-License-Identifier: MPL-2.0

package agent

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAgentModuleSupply_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + testdata.AgentCreateConfig + AgentModuleSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_module_supply.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceAgentModuleSupply_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + testdata.AgentCreateConfig + AgentModuleSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_module_supply.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_agent_module_supply.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
