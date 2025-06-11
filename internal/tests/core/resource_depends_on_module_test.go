// SPDX-License-Identifier: MPL-2.0

package core

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)


func TestAccResourceDependsOnModule_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleCreateConfig + ModuleCreateConfigDeltaTwo + DependsOnModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "id"),
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "module_id"),
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "depends_on_module_id"),
				),
			},
		},
	})
}

func TestAccResourceDependsOnModule_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleCreateConfig + ModuleCreateConfigDeltaTwo + DependsOnModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "id"),
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "module_id"),
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "depends_on_module_id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + ModuleCreateConfig + ModuleCreateConfigDeltaTwo + ModuleCreateConfigDeltaThree + `
resource "snapcd_depends_on_module" "this" { 
  module_id = snapcd_module.this.id
  depends_on_module_id = snapcd_module.three.id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "id"),
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "module_id"),
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "depends_on_module_id"),
				),
			},
		},
	})
}

func TestAccResourceDependsOnModule_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleCreateConfig + ModuleCreateConfigDeltaTwo + DependsOnModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_depends_on_module.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_depends_on_module.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
