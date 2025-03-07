// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceEnvVarFromSecretCreateConfig = secret.AzureKeyVaultSecretScopedToNamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_env_var_from_secret" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	   = "somevalue%s"
  secret_name  = snapcd_azure_key_vault_secret_scoped_to_namespace.this.name
  secret_scope        = "Namespace"
}
  
`)

func TestAccResourceNamespaceEnvVarFromSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceEnvVarFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceEnvVarFromSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceEnvVarFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_env_var_from_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_env_var_from_secret" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  = "somevalue%s"
  secret_value  = "barrr"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_env_var_from_secret.this", "secret_value", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceEnvVarFromSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceEnvVarFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_env_var_from_secret.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
