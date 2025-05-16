// SPDX-License-Identifier: MPL-2.0

package secret_store_assignment

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret_store"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecretStoreNamespaceAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + AzureKeyVaultSecretStoreNamespaceAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_namespace_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceSecretStoreNamespaceAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + AzureKeyVaultSecretStoreNamespaceAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_namespace_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_secret_store_namespace_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
