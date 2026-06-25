// SPDX-License-Identifier: MPL-2.0

package agent

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAgent_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.AgentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent.this", "id"),
					resource.TestCheckResourceAttr("snapcd_agent.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}

func TestAccResourceAgent_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.AgentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent.this", "id"),
					resource.TestCheckResourceAttr("snapcd_agent.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.AgentServicePrincipalConfig + providerconfig.AppendRandomString(`
resource "snapcd_agent" "this" {
  name                 = "someNEWvalue%s"
  service_principal_id = data.snapcd_service_principal.agent.id
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent.this", "id"),
					resource.TestCheckResourceAttr("snapcd_agent.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceAgent_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.AgentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_agent.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_agent.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
