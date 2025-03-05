// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceParamFromLiteralCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_param_from_literal" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	= "somevalue%s"
  literal_value  	= "bar"
}
  
`)

func TestAccResourceNamespaceParamFromLiteral_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceParamFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_literal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceParamFromLiteral_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceParamFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_param_from_literal.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + `
resource "snapcd_namespace_param_from_literal" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  = "somevalue%s"
  literal_value  = "barrr"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_param_from_literal.this", "literal_value", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceParamFromLiteral_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceParamFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_literal.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_param_from_literal.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
