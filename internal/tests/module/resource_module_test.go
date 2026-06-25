// SPDX-License-Identifier: MPL-2.0

package module

import (
	"strings"
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceModule_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + ModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModule_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + ModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module.this", "apply_approval_threshold", "1"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + strings.Replace(ModuleCreateConfig, "apply_approval_threshold               = 1", "apply_approval_threshold               = 2", -1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module.this", "apply_approval_threshold", "2"),
				),
			},
		},
	})
}

func TestAccResourceModule_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + ModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
