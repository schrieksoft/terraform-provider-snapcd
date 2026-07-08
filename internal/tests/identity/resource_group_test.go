// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceGroup_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + GroupCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceGroup_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + GroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group.this", "id"),
					resource.TestCheckResourceAttr("snapcd_group.this", "name", providerconfig.AppendRandomString("grp_update_%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + GroupUpdateNewConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group.this", "id"),
					resource.TestCheckResourceAttr("snapcd_group.this", "name", providerconfig.AppendRandomString("grp_update_new_%s")),
				),
			},
		},
	})
}

func TestAccResourceGroup_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + GroupImportConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_group.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_group.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
