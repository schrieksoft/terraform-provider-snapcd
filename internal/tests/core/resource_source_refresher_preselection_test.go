// SPDX-License-Identifier: MPL-2.0

package core

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)



func TestAccResourceSourceRefresherPreselection_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + RunnerPoolCreateConfig + SourceRefresherPreselectionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_source_refresher_preselection.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceSourceRefresherPreselection_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + RunnerPoolCreateConfig + SourceRefresherPreselectionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_source_refresher_preselection.this", "id"),
					resource.TestCheckResourceAttr("snapcd_source_refresher_preselection.this", "source_url", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + RunnerPoolCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_source_refresher_preselection" "this" { 
  source_url = "someNEWvalue%s"
  runner_pool_id = snapcd_runner_pool.this.id
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_source_refresher_preselection.this", "id"),
					resource.TestCheckResourceAttr("snapcd_source_refresher_preselection.this", "source_url", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceSourceRefresherPreselection_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + RunnerPoolCreateConfig + SourceRefresherPreselectionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_source_refresher_preselection.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_source_refresher_preselection.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
