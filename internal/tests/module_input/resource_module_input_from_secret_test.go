// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"terraform-provider-snapcd/internal/tests/secret_store"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleInputFromSecretCreateConfig = secret_store.AwsSecretsManagerSecretStoreCreateConfig + secret.SecretScopedToStackCreateConfigDelta + providerconfig.AppendRandomString(`
resource "snapcd_module_input_from_secret" "this" { 
  input_kind = "Param"
  module_id  = snapcd_module.this.id
  name  	 = "somevalue%s"
  secret_id  = snapcd_secret_scoped_to_stack.this.id
  type 		 = "NotString"
}
  
`)

func TestAccResourceModuleInputFromSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleInputFromSecretCreateConfig,
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
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleInputFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_input_from_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + secret_store.AwsSecretsManagerSecretStoreCreateConfig + secret.SecretScopedToStackCreateConfigDelta + providerconfig.AppendRandomString(`
resource "snapcd_module_input_from_secret" "this" { 
  input_kind = "Param"
  module_id  = snapcd_module.this.id
  name  	 = "somevalue%s"
  secret_id  = snapcd_simple_secret_scoped_to_stack.this.id
  type 		 = "NotString"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_input_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_input_from_secret.this", "type", "NotString"),
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
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + ModuleInputFromSecretCreateConfig,
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
