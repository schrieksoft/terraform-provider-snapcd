// SPDX-License-Identifier: MPL-2.0

package secret_store

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret_store"
	"terraform-provider-snapcd/internal/tests/core"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)


func TestAccResourceSecretStoreModuleAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + SecretStoreModuleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_module_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceSecretStoreModuleAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig+ SecretStoreModuleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_module_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_secret_store_module_assignment.this", "permission", "ReadWrite"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_secret_store_module_assignment" "this" { 
  permission  		= "Read" 
  secret_store_id   = snapcd_secret_store_module_assignment.this.id
  secret_store_id   = snapcd_namespace.this.id
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_module_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_secret_store_module_assignment.this", "permission", "Read"),
				),
			},
		},
	})
}

func TestAccResourceSecretStoreModuleAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig+ SecretStoreModuleAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_module_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_secret_store_module_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
