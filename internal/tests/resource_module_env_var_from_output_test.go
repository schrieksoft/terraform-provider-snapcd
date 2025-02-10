// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var moduleEnvVarFromOutputCreateConfig = appendRandomString(`
resource "snapcd_module_env_var_from_output" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  module_name  	= "bar"
  namespace_name  	= "bar"
  output_name   = "bar"
}
  
`)

func TestAccResourceModuleEnvVarFromOutput_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromOutputCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_output.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromOutput_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromOutputCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_output.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_output.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + moduleCreateConfig + appendRandomString(`
resource "snapcd_module_env_var_from_output" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  module_name  	= "bar"
  namespace_name  	= "barrr"
  output_name   = "bar"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_output.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_output.this", "namespace_name", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromOutput_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromOutputCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_output.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_env_var_from_output.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
