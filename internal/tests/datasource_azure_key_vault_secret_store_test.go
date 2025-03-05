// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceAzureKeyVaultSecretStore(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + azureKeyVaultSecretStoreCreateConfig + `

data "snapcd_azure_key_vault_secret_store" "this" {
	name = snapcd_azure_key_vault_secret_store.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_azure_key_vault_secret_store.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_azure_key_vault_secret_store.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
