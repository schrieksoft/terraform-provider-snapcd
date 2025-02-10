// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var moduleParamFromOutputSetCreateConfig = appendRandomString(`
resource "snapcd_module_param_from_output_set" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  module_name  	= "bar"
  namespace_name  	= "bar"
}
  
`)

func TestAccResourceModuleParamFromOutputSet_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromOutputSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_output_set.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromOutputSet_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromOutputSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_output_set.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_output_set.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + moduleCreateConfig + appendRandomString(`
resource "snapcd_module_param_from_output_set" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  module_name  	= "bar"
  namespace_name  	= "barrr"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_output_set.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_output_set.this", "namespace_name", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromOutputSet_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromOutputSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_output_set.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_param_from_output_set.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
