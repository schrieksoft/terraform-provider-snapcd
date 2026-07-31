// SPDX-License-Identifier: MPL-2.0

package trigger_paths

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceAdditionalTriggerPathCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_additional_trigger_path" "this" {
  namespace_id = snapcd_namespace.this.id
  path         = "shared/config%s"
}

`)

func TestAccResourceNamespaceAdditionalTriggerPath_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceAdditionalTriggerPathCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_additional_trigger_path.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceAdditionalTriggerPath_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceAdditionalTriggerPathCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_additional_trigger_path.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_additional_trigger_path.this", "path", providerconfig.AppendRandomString("shared/config%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_additional_trigger_path" "this" {
  namespace_id = snapcd_namespace.this.id
  path         = "shared/other%s"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_additional_trigger_path.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_additional_trigger_path.this", "path", providerconfig.AppendRandomString("shared/other%s")),
				),
			},
		},
	})
}

func TestAccResourceNamespaceAdditionalTriggerPath_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceAdditionalTriggerPathCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_additional_trigger_path.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_additional_trigger_path.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
