// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleParamFromLiteralCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_param_from_literal" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  literal_value  	= "bar"
}
  
`)

func TestAccResourceModuleParamFromLiteral_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_literal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromLiteral_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_literal.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_param_from_literal" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  literal_value  = "barrr"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_literal.this", "literal_value", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromLiteral_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_literal.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_param_from_literal.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
