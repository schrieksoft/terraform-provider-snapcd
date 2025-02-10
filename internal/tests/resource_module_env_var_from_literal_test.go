// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var moduleEnvVarFromLiteralCreateConfig = appendRandomString(`
resource "snapcd_module_env_var_from_literal" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  literal_value  	= "bar"
}
  
`)

func TestAccResourceModuleEnvVarFromLiteral_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_literal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromLiteral_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_literal.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + moduleCreateConfig + appendRandomString(`
resource "snapcd_module_env_var_from_literal" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  literal_value  = "barrr"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_literal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_literal.this", "literal_value", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromLiteral_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromLiteralCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_literal.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_env_var_from_literal.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
