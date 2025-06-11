// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleParamFromSecretScopedToStackCreateConfig = secret.SimpleSecretScopedToStackCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_param_from_secret_scoped_to_stack" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  secret_scoped_to_stack_id = snapcd_simple_secret_scoped_to_stack.this.id
}
  
`)

func TestAccResourceModuleParamFromSecretScopedToStack_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromSecretScopedToStackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret_scoped_to_stack.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromSecretScopedToStack_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromSecretScopedToStackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret_scoped_to_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_secret_scoped_to_stack.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + secret.SimpleSecretScopedToStackCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_param_from_secret_scoped_to_stack" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  secret_scoped_to_stack_id = snapcd_simple_secret_scoped_to_stack.this.id
  type = "NotString"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret_scoped_to_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_secret_scoped_to_stack.this", "type", "NotString"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromSecretScopedToStack_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleParamFromSecretScopedToStackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret_scoped_to_stack.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_param_from_secret_scoped_to_stack.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
