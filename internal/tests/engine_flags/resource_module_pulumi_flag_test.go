// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModulePulumiFlagCreateConfig = `
resource "snapcd_module_pulumi_flag" "this" {
  module_id 	= snapcd_module.this.id
  task  		= "Init"
  flag  		= "Refresh"
  value  		= "true"
}

`

func TestAccResourceModulePulumiFlag_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModulePulumiFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_flag.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModulePulumiFlag_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModulePulumiFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_flag.this", "value", "true"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_flag.this", "task", "Init"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_flag.this", "flag", "Refresh"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + `
resource "snapcd_module_pulumi_flag" "this" {
  module_id  	= snapcd_module.this.id
  task  		= "Init"
  flag  		= "Refresh"
  value   	= "false"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_flag.this", "value", "false"),
				),
			},
		},
	})
}

func TestAccResourceModulePulumiFlag_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModulePulumiFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_flag.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_pulumi_flag.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
