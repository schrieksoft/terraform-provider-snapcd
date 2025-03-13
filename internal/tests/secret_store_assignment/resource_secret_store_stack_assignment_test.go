// SPDX-License-Identifier: MPL-2.0

package secret_store_assignment

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret_store"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecretStoreStackAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + AzureKeyVaultSecretStoreStackAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_stack_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceSecretStoreStackAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + AzureKeyVaultSecretStoreStackAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_stack_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_secret_store_stack_assignment.this", "permission", "ReadWrite"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_secret_store_stack_assignment" "this" { 
  secret_store_id   = snapcd_secret_store_stack_assignment.this.id
  stack_id   = snapcd_stack.this.id
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_stack_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_secret_store_stack_assignment.this", "permission", "Read"),
				),
			},
		},
	})
}

func TestAccResourceSecretStoreStackAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + secret_store.AzureKeyVaultSecretStoreCreateConfig + AzureKeyVaultSecretStoreStackAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_store_stack_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_secret_store_stack_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
