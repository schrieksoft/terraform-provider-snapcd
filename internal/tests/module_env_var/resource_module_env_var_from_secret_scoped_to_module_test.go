// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleEnvVarFromSecretScopedToModuleCreateConfig = secret.SimpleSecretScopedToModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_env_var_from_secret_scoped_to_module" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  secret_scoped_to_module_id = snapcd_simple_secret_scoped_to_module.this.id
}
  
`)

func TestAccResourceModuleEnvVarFromSecretScopedToModule_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleEnvVarFromSecretScopedToModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_secret_scoped_to_module.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromSecretScopedToModule_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleEnvVarFromSecretScopedToModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_secret_scoped_to_module.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_secret_scoped_to_module.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + secret.SimpleSecretScopedToModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_env_var_from_secret_scoped_to_module" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  secret_scoped_to_module_id = snapcd_simple_secret_scoped_to_module.this.id
  type = "NotString"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_secret_scoped_to_module.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_secret_scoped_to_module.this", "type", "NotString"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromSecretScopedToModule_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleEnvVarFromSecretScopedToModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_secret_scoped_to_module.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_env_var_from_secret_scoped_to_module.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
