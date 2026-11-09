// SPDX-License-Identifier: MPL-2.0

package output

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceModuleOutputDownstreamRestriction_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleOutputDownstreamRestrictionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_output_downstream_restriction.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleOutputDownstreamRestriction_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleOutputDownstreamRestrictionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_output_downstream_restriction.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_output_downstream_restriction.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
