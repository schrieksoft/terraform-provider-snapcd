// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var namespaceParamFromLiteralCreateConfig = appendRandomString(`
resource "snapcd_namespace_param_from_literal" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	= "somevalue%s"
  literal_value  	= "bar"
}
  
`)

func TestAccResourceNamespaceParamFromLiteral_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceParamFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_literal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceParamFromLiteral_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceParamFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_param_from_literal.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + namespaceCreateConfig + `
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceParamFromLiteralCreateConfig,
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
