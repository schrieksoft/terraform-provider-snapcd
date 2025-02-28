// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceEnvVarFromDefinition(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig + namespaceEnvVarFromDefinitionCreateConfig + `
data "snapcd_namespace_env_var_from_definition" "this" {
	name 		= snapcd_namespace_env_var_from_definition.this.name
	namespace_id 	= snapcd_namespace_env_var_from_definition.this.namespace_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_env_var_from_definition.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_env_var_from_definition.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
