// SPDX-License-Identifier: MPL-2.0

package policies

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespacePulumiInlinePolicyCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_pulumi_inline_policy" "this" {
  name      = "mypolicy%s"
  namespace_id = snapcd_namespace.this.id

  runtime        = "Python"
  policy_content = "# CrossGuard pack entry module placeholder"
}

`)

func TestAccResourceNamespacePulumiInlinePolicy_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespacePulumiInlinePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_inline_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_inline_policy.this", "enabled", "true"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_inline_policy.this", "evaluate_on", "ApplyOnly"),
				),
			},
		},
	})
}

func TestAccResourceNamespacePulumiInlinePolicy_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespacePulumiInlinePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_inline_policy.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_pulumi_inline_policy" "this" {
  name        = "mypolicy%s"
  namespace_id   = snapcd_namespace.this.id
  enabled     = false

  runtime        = "Python"
  policy_content = "# CrossGuard pack entry module placeholder"
}

`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_inline_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_inline_policy.this", "enabled", "false"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_inline_policy.this", "evaluate_on", "ApplyOnly"),
				),
			},
		},
	})
}

func TestAccResourceNamespacePulumiInlinePolicy_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespacePulumiInlinePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_inline_policy.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_pulumi_inline_policy.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
