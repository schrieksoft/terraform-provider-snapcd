// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var azureKeyVaultSecretStoreCreateConfig = appendRandomString(`
resource "snapcd_azure_key_vault_secret_store" "this" { 
  name  = "somevalue%s"
  key_vault_url = "https://snapcdlocaltesting.vault.azure.net/"
}`)

func TestAccResourceAzureKeyVaultSecretStore_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + azureKeyVaultSecretStoreCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_azure_key_vault_secret_store.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceAzureKeyVaultSecretStore_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + azureKeyVaultSecretStoreCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_azure_key_vault_secret_store.this", "id"),
					resource.TestCheckResourceAttr("snapcd_azure_key_vault_secret_store.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + appendRandomString(`
resource "snapcd_azure_key_vault_secret_store" "this" { 
  name = "someNEWvalue%s"
  key_vault_url = "https://snapcdlocaltesting.vault.azure.net/"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_azure_key_vault_secret_store.this", "id"),
					resource.TestCheckResourceAttr("snapcd_azure_key_vault_secret_store.this", "name", appendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceAzureKeyVaultSecretStore_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + azureKeyVaultSecretStoreCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_azure_key_vault_secret_store.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_azure_key_vault_secret_store.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
