// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleInputFromNamespaceCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_input_from_namespace" "this" { 
  input_kind 		= "Param"
  module_id 		= snapcd_module.this.id
  name  			= "somevalue%s"
  reference_name  	= "bar"
}
  
`)

func TestAccResourceModuleInputFromNamespace_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleInputFromNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_namespace.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleInputFromNamespace_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleInputFromNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_input_from_namespace.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_input_from_namespace" "this" { 
  input_kind 	  = "Param"
  module_id 	  = snapcd_module.this.id
  name  		  = "somevalue%s"
  reference_name  = "barrr"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_input_from_namespace.this", "reference_name", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceModuleInputFromNamespace_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleInputFromNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_namespace.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_input_from_namespace.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
