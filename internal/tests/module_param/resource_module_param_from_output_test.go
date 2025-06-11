// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleParamFromOutputCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_param_from_output" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  output_module_id  	= snapcd_module.two.id
  output_name   = "bar"
}
  
`)

func TestAccResourceModuleParamFromOutput_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromOutputCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_output.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromOutput_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromOutputCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_output.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_output.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_param_from_output" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  output_module_id  	= snapcd_module.this.id
  output_name   = "bar"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_output.this", "id"),
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_output.this", "output_module_id"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromOutput_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromOutputCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_output.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_param_from_output.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
