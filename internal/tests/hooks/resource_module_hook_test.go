// SPDX-License-Identifier: MPL-2.0

package hooks

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleHookCreateConfig = `
resource "snapcd_module_hook" "this" {
  module_id 	= snapcd_module.this.id
  task  		= "Plan"
  phase 		= "Before"
  script  		= "echo 'before plan'"
}

`

func TestAccResourceModuleHook_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleHookCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_hook.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_hook.this", "task", "Plan"),
					resource.TestCheckResourceAttr("snapcd_module_hook.this", "phase", "Before"),
					resource.TestCheckResourceAttr("snapcd_module_hook.this", "script", "echo 'before plan'"),
				),
			},
		},
	})
}

func TestAccResourceModuleHook_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleHookCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_hook.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_hook.this", "task", "Plan"),
					resource.TestCheckResourceAttr("snapcd_module_hook.this", "phase", "Before"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + `
resource "snapcd_module_hook" "this" {
  module_id 	= snapcd_module.this.id
  task  		= "Plan"
  phase 		= "Before"
  script  		= "echo 'updated script'"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_hook.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_hook.this", "script", "echo 'updated script'"),
				),
			},
		},
	})
}

func TestAccResourceModuleHook_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleHookCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_hook.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_hook.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
