// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleEnvVarFromDefinitionCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_env_var_from_definition" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  definition_name  	= "ModuleName"
}
  
`)

func TestAccResourceModuleEnvVarFromDefinition_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleEnvVarFromDefinitionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_definition.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromDefinition_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleEnvVarFromDefinitionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_definition.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_definition.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_env_var_from_definition" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  definition_name  = "NamespaceName"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_definition.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_definition.this", "definition_name", "NamespaceName"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromDefinition_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleEnvVarFromDefinitionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_definition.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_env_var_from_definition.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
