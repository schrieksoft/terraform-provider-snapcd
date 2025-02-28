// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var moduleParamFromNamespaceCreateConfig = appendRandomString(`
resource "snapcd_module_param_from_namespace" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  reference_name  	= "bar"
}
  
`)

func TestAccResourceModuleParamFromNamespace_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_namespace.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromNamespace_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_namespace.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + moduleCreateConfig + appendRandomString(`
resource "snapcd_module_param_from_namespace" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  reference_name  = "barrr"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_namespace.this", "reference_name", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromNamespace_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_namespace.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_param_from_namespace.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
