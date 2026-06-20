// SPDX-License-Identifier: MPL-2.0

package agent_supply

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAgentNamespaceSupply_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + core.AgentCreateConfig + AgentNamespaceSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_namespace_supply.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceAgentNamespaceSupply_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + core.AgentCreateConfig + AgentNamespaceSupplyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent_namespace_supply.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_agent_namespace_supply.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
