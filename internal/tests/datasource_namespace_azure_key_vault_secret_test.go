// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceAzureKeyVaultSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceAzureKeyVaultSecretCreateConfig + `

data "snapcd_namespace_azure_key_vault_secret" "this" {
	name = snapcd_namespace_azure_key_vault_secret.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_azure_key_vault_secret.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_azure_key_vault_secret.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
