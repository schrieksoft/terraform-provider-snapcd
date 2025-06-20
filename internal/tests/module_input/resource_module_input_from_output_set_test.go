// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleInputFromOutputSetCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_input_from_output_set" "this" { 
  input_kind 		= "Param"
  module_id 		= snapcd_module.this.id
  name  			= "somevalue%s"
  output_module_id  = snapcd_module.two.id
}
  
`)

func TestAccResourceModuleInputFromOutputSet_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + core.ModuleCreateConfigDeltaTwo + ModuleInputFromOutputSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_output_set.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleInputFromOutputSet_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + core.ModuleCreateConfigDeltaTwo + ModuleInputFromOutputSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_output_set.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_input_from_output_set.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + core.ModuleCreateConfigDeltaTwo + providerconfig.AppendRandomString(`
resource "snapcd_module_input_from_output_set" "this" { 
  input_kind 		= "Param"
  module_id 		= snapcd_module.this.id
  name  			= "somevalue%s"
  output_module_id  = snapcd_module.two.id
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_output_set.this", "id"),
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_output_set.this", "output_module_id"),
				),
			},
		},
	})
}

func TestAccResourceModuleInputFromOutputSet_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + core.ModuleCreateConfigDeltaTwo + ModuleInputFromOutputSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_output_set.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_input_from_output_set.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
