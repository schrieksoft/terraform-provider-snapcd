// SPDX-License-Identifier: MPL-2.0

package secret_store

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAzureSecretStore_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + AzureSecretStoreCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_azure_secret_store.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceAzureSecretStore_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + AzureSecretStoreCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_azure_secret_store.this", "id"),
					resource.TestCheckResourceAttr("snapcd_azure_secret_store.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_azure_secret_store" "this" { 
  name = "someNEWvalue%s"
  key_vault_url = "https://snapcdlocaltesting.vault.azure.net/"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_azure_secret_store.this", "id"),
					resource.TestCheckResourceAttr("snapcd_azure_secret_store.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceAzureSecretStore_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + AzureSecretStoreCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_azure_secret_store.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_azure_secret_store.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
