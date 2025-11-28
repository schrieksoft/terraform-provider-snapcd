// SPDX-License-Identifier: MPL-2.0

package core

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSourceRefresherPreselection(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + RunnerCreateConfig + SourceRefresherPreselectionCreateConfig + `

data "snapcd_source_refresher_preselection" "this" {
	source_url = snapcd_source_refresher_preselection.this.source_url
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_source_refresher_preselection.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_source_refresher_preselection.this", "source_url", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
