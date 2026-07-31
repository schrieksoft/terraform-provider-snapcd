// SPDX-License-Identifier: MPL-2.0

package trigger_paths

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleAdditionalTriggerPathCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_additional_trigger_path" "this" {
  module_id = snapcd_module.this.id
  path      = "shared/scripts%s"
}

`)

func TestAccResourceModuleAdditionalTriggerPath_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleAdditionalTriggerPathCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_additional_trigger_path.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleAdditionalTriggerPath_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleAdditionalTriggerPathCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_additional_trigger_path.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_additional_trigger_path.this", "path", providerconfig.AppendRandomString("shared/scripts%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_additional_trigger_path" "this" {
  module_id = snapcd_module.this.id
  path      = "shared/other%s"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_additional_trigger_path.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_additional_trigger_path.this", "path", providerconfig.AppendRandomString("shared/other%s")),
				),
			},
		},
	})
}

func TestAccResourceModuleAdditionalTriggerPath_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleAdditionalTriggerPathCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_additional_trigger_path.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_additional_trigger_path.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
