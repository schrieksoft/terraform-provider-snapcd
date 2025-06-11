// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceEnvVarFromSecretScopedToNamespaceCreateConfig = secret.SimpleSecretScopedToNamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_env_var_from_secret_scoped_to_namespace" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	= "somevalue%s"
  secret_scoped_to_namespace_id = snapcd_simple_secret_scoped_to_namespace.this.id
}
  
`)

func TestAccResourceNamespaceEnvVarFromSecretScopedToNamespace_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceEnvVarFromSecretScopedToNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret_scoped_to_namespace.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceEnvVarFromSecretScopedToNamespace_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceEnvVarFromSecretScopedToNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret_scoped_to_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_env_var_from_secret_scoped_to_namespace.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + secret.SimpleSecretScopedToNamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_env_var_from_secret_scoped_to_namespace" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  = "somevalue%s"
  secret_scoped_to_namespace_id = snapcd_simple_secret_scoped_to_namespace.this.id
  usage_mode = "UseByDefault"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret_scoped_to_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_env_var_from_secret_scoped_to_namespace.this", "usage_mode", "UseByDefault"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceEnvVarFromSecretScopedToNamespace_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceEnvVarFromSecretScopedToNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret_scoped_to_namespace.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_env_var_from_secret_scoped_to_namespace.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
