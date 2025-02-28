// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var namespaceParamFromDefinitionCreateConfig = appendRandomString(`
resource "snapcd_namespace_param_from_definition" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	= "somevalue%s"
  definition_name  	= "ModuleName"
}
  
`)

func TestAccResourceNamespaceParamFromDefinition_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceParamFromDefinitionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_definition.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceParamFromDefinition_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceParamFromDefinitionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_definition.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_param_from_definition.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + namespaceCreateConfig + `
resource "snapcd_namespace_param_from_definition" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  = "somevalue%s"
  definition_name  = "NamespaceName"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_definition.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_param_from_definition.this", "definition_name", "NamespaceName"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceParamFromDefinition_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceParamFromDefinitionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_definition.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_param_from_definition.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
