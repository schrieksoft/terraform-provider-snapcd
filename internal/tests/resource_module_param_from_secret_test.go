// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var moduleParamFromSecretCreateConfig = appendRandomString(`
resource "snapcd_module_param_from_secret" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  secret_value  	= "bar"
}
  
`)

func TestAccResourceModuleParamFromSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_secret.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + moduleCreateConfig + appendRandomString(`
resource "snapcd_module_param_from_secret" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  secret_value  = "barrr"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_secret.this", "secret_value", "barrr"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig + moduleParamFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_param_from_secret.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
