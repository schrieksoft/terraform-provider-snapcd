// SPDX-License-Identifier: MPL-2.0

package engine_flags

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
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
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + ModuleTerraformFlagCreateConfig,
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
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + ModuleTerraformFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_terraform_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_flag.this", "task", "Init"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_flag.this", "flag", "Upgrade"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + `
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
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + ModuleTerraformFlagCreateConfig,
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
