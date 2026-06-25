// SPDX-License-Identifier: MPL-2.0

package agent

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAgentStackSupply_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.StackCreateConfig + testdata.AgentCreateConfig + AgentStackSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_stack_supply.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceAgentStackSupply_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.StackCreateConfig + testdata.AgentCreateConfig + AgentStackSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_stack_supply.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_agent_stack_supply.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
