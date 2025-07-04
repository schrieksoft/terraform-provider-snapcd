// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"terraform-provider-snapcd/internal/tests/secret_store"
	"terraform-provider-snapcd/internal/tests/secret_store_assignment"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceInputFromSecretCreateConfig = secret_store.AwsSecretStoreCreateConfig + secret_store_assignment.AwsSecretStoreStackAssignmentCreateConfig + secret.SecretScopedToStackCreateConfigDelta + providerconfig.AppendRandomString(`
resource "snapcd_namespace_input_from_secret" "this" { 
  input_kind 	= "Param"
  namespace_id  = snapcd_namespace.this.id
  name  		= "somevalue%s"
  secret_id 	= snapcd_secret_scoped_to_stack.this.id
}
`)

var NamespaceInputFromSecretCreateConfigNew = secret_store.AwsSecretStoreCreateConfig + secret_store_assignment.AwsSecretStoreStackAssignmentCreateConfig + secret.SecretScopedToStackCreateConfigDelta + providerconfig.AppendRandomString(`
resource "snapcd_namespace_input_from_secret" "this" { 
  input_kind 	= "Param"
  namespace_id  = snapcd_namespace.this.id
  name  		= "someNEWvalue%s"
  secret_id 	= snapcd_secret_scoped_to_stack.this.id
}
  
`)

func TestAccResourceNamespaceInputFromSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceInputFromSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_input_from_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfigNew,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_input_from_secret.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceNamespaceInputFromSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_secret.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_input_from_secret.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
