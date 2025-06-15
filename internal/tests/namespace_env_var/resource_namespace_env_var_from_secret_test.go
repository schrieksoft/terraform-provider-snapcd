// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceEnvVarFromSecretCreateConfig = secret.SimpleSecretCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_env_var_from_secret" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	= "somevalue%s"
  secret_id = snapcd_simple_secret.this.id
}
  
`)

func TestAccResourceNamespaceEnvVarFromSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceEnvVarFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceEnvVarFromSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceEnvVarFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_env_var_from_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + secret.SimpleSecretCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_env_var_from_secret" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  = "somevalue%s"
  secret_id = snapcd_simple_secret.this.id
  usage_mode = "UseByDefault"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_env_var_from_secret.this", "usage_mode", "UseByDefault"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceEnvVarFromSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceEnvVarFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_secret.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_env_var_from_secret.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
