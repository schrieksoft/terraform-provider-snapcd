// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleExtraFileCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_extra_file" "this" { 
  module_id 	= snapcd_module.this.id
  file_name  	= "somevalue%s"
  contents  	= "foo"
}
  
`)

func TestAccResourceModuleExtraFile_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleExtraFileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_extra_file.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleExtraFile_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleExtraFileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_extra_file.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_extra_file.this", "contents", "foo"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_extra_file" "this" { 
  module_id  = snapcd_module.this.id
  file_name  = "somevalue%s"
  contents   = "bar"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_extra_file.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_extra_file.this", "contents", "bar"),
				),
			},
		},
	})
}

func TestAccResourceModuleExtraFile_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleExtraFileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_extra_file.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_extra_file.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
