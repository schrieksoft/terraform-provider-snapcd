// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceInputFromDefinitionCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_input_from_definition" "this" { 
  input_kind 	  = "Param"
  namespace_id 	  = snapcd_namespace.this.id
  name  		  = "somevalue%s"
  definition_name = "ModuleName"
}
  
`)

func TestAccResourceNamespaceInputFromDefinition_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromDefinitionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_definition.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceInputFromDefinition_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromDefinitionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_definition.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_input_from_definition.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + `
resource "snapcd_namespace_input_from_definition" "this" { 
  input_kind 	  = "Param"
  namespace_id    = snapcd_namespace.this.id
  name  		  = "somevalue%s"
  definition_name = "NamespaceName"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_definition.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_input_from_definition.this", "definition_name", "NamespaceName"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceInputFromDefinition_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromDefinitionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_definition.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_input_from_definition.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
