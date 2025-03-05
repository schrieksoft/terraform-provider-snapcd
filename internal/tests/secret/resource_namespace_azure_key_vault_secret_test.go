// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceAzureKeyVaultSecretCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_azure_key_vault_secret" "this" { 
  name  = "somevalue%s"
}`)

func TestAccResourceNamespaceAzureKeyVaultSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceAzureKeyVaultSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_azure_key_vault_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceAzureKeyVaultSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceAzureKeyVaultSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_azure_key_vault_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_azure_key_vault_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_azure_key_vault_secret" "this" { 
  name = "someNEWvalue%s"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_azure_key_vault_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_azure_key_vault_secret.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceNamespaceAzureKeyVaultSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceAzureKeyVaultSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_azure_key_vault_secret.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_azure_key_vault_secret.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
