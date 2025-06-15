// SPDX-License-Identifier: MPL-2.0

package core

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var DependsOnModuleDataSourceConfig = providerconfig.AppendRandomString(`
data "snapcd_depends_on_module" "this" {
  id = snapcd_depends_on_module.this.id
}
`)

func TestAccDataSourceDependsOnModule_Read(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ModuleCreateConfig + ModuleCreateConfigDeltaTwo + DependsOnModuleCreateConfig + DependsOnModuleDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_depends_on_module.this", "id"),
					resource.TestCheckResourceAttrSet("data.snapcd_depends_on_module.this", "module_id"),
					resource.TestCheckResourceAttrSet("data.snapcd_depends_on_module.this", "depends_on_module_id"),
				),
			},
		},
	})
}
