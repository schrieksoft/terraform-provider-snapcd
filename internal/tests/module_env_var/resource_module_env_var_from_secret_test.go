// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleEnvVarFromSecretCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_env_var_from_secret" "this" {
  name  	   = "somevalue%s"
  module_id    = snapcd_module.this.id
  secret_name  = "name-in-db"
  secret_scope = "Namespace"
}  
`)

func TestAccResourceModuleEnvVarFromSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleEnvVarFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleEnvVarFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_env_var_from_secret" "this" { 
  name  	   = "somevalue%s"
  module_id    = snapcd_module.this.id
  secret_name  = "name-in-db-NEW"
  secret_scope = "Namespace"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_secret.this", "secret_name", "name-in-db-NEW"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleEnvVarFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_secret.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_env_var_from_secret.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
