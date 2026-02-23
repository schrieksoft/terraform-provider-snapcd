// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleTerraformFlagCreateConfig = `
resource "snapcd_module_terraform_flag" "this" {
  module_id 	= snapcd_module.this.id
  task  		= "Init"
  flag  		= "Upgrade"
}

`

func TestAccResourceModuleTerraformFlag_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleTerraformFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_terraform_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_flag.this", "task", "Init"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_flag.this", "flag", "Upgrade"),
				),
			},
		},
	})
}

func TestAccResourceModuleTerraformFlag_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleTerraformFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_terraform_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_flag.this", "task", "Init"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_flag.this", "flag", "Upgrade"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + `
resource "snapcd_module_terraform_flag" "this" {
  module_id  	= snapcd_module.this.id
  task  		= "Init"
  flag  		= "Upgrade"
  value   	= "true"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_terraform_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_flag.this", "value", "true"),
				),
			},
		},
	})
}

func TestAccResourceModuleTerraformFlag_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleTerraformFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_terraform_flag.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_terraform_flag.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
