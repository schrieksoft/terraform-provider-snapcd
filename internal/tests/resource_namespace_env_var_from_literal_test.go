// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var namespaceEnvVarFromLiteralCreateConfig = appendRandomString(`
resource "snapcd_namespace_env_var_from_literal" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	= "somevalue%s"
  literal_value  	= "bar"
}
  
`)

func TestAccResourceNamespaceEnvVarFromLiteral_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceEnvVarFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_literal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceEnvVarFromLiteral_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceEnvVarFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_env_var_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_env_var_from_literal.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + namespaceCreateConfig + appendRandomString(`
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceEnvVarFromLiteralCreateConfig,
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
