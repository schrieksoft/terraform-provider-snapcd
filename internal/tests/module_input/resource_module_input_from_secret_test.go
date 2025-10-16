// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleInputFromSecretCreateConfig = core.StackDebugDataSourceDelta + core.NamespaceDebugDataSourceDelta + core.ModuleDebugDataSourceDelta + core.StackSecretDebugDataSourceDelta + providerconfig.AppendRandomString(`

resource "snapcd_module_input_from_secret" "this" { 
  input_kind 	= "Param"
  module_id  	= data.snapcd_module.debug.id
  name  		= "somevalue%s"
  secret_id 	= data.snapcd_stack_secret.debug.id
}
`)

var ModuleInputFromSecretCreateConfigNew = core.StackDebugDataSourceDelta + core.NamespaceDebugDataSourceDelta + core.ModuleDebugDataSourceDelta + core.StackSecretDebugDataSourceDelta + providerconfig.AppendRandomString(`

resource "snapcd_module_input_from_secret" "this" { 
  input_kind 	= "Param"
  module_id  	= data.snapcd_module.debug.id
  name  		= "someNEWvalue%s"
  secret_id 	= data.snapcd_stack_secret.debug.id
}
`)

func TestAccResourceModuleInputFromSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleInputFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleInputFromSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleInputFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_input_from_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + ModuleInputFromSecretCreateConfigNew,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_input_from_secret.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceModuleInputFromSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleInputFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_secret.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_input_from_secret.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
