// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceEnvVarFromLiteralCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_env_var_from_literal" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	= "somevalue%s"
  literal_value  	= "bar"
}
  
`)

func TestAccResourceNamespaceEnvVarFromLiteral_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceEnvVarFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_literal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceEnvVarFromLiteral_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceEnvVarFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_env_var_from_literal.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_env_var_from_literal" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  = "somevalue%s"
  literal_value  = "barrr"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_env_var_from_literal.this", "literal_value", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceEnvVarFromLiteral_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceEnvVarFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_literal.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_env_var_from_literal.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
