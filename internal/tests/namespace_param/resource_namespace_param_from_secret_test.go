// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceParamFromSecretCreateConfig = secret.AzureKeyVaultSecretScopedToNamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_param_from_secret" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	   = "somevalue%s"
  secret_name  = snapcd_azure_key_vault_secret_scoped_to_namespace.this.name
  scope        = "Namespace"
}
  
`)

func TestAccResourceNamespaceParamFromSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceParamFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceParamFromSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceParamFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_param_from_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + `
resource "snapcd_namespace_param_from_secret" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  = "somevalue%s"
  secret_value  = "barrr"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_param_from_secret.this", "secret_value", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceParamFromSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceParamFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_secret.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_param_from_secret.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
