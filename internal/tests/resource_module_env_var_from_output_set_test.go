// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var moduleEnvVarFromOutputSetCreateConfig = appendRandomString(`
resource "snapcd_module_env_var_from_output_set" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  module_name  	= "bar"
  namespace_name  	= "bar"
}
  
`)

func TestAccResourceModuleEnvVarFromOutputSet_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromOutputSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_output_set.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromOutputSet_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromOutputSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_output_set.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_output_set.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + moduleCreateConfig + appendRandomString(`
resource "snapcd_module_env_var_from_output_set" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  module_name  	= "bar"
  namespace_name  	= "barrr"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_output_set.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_env_var_from_output_set.this", "namespace_name", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceModuleEnvVarFromOutputSet_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleEnvVarFromOutputSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_env_var_from_output_set.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_env_var_from_output_set.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
