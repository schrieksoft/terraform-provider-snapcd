// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var namespaceAzureKeyVaultSecretCreateConfig = appendRandomString(`
resource "snapcd_namespace_azure_key_vault_secret" "this" { 
  name  = "somevalue%s"
}`)

func TestAccResourceNamespaceAzureKeyVaultSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceAzureKeyVaultSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_azure_key_vault_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceAzureKeyVaultSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceAzureKeyVaultSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_azure_key_vault_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_azure_key_vault_secret.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + appendRandomString(`
resource "snapcd_namespace_azure_key_vault_secret" "this" { 
  name = "someNEWvalue%s"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_azure_key_vault_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_azure_key_vault_secret.this", "name", appendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceNamespaceAzureKeyVaultSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceAzureKeyVaultSecretCreateConfig,
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
