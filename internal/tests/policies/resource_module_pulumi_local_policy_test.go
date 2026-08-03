// SPDX-License-Identifier: MPL-2.0

package policies

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModulePulumiLocalPolicyCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_pulumi_local_policy" "this" {
  name      = "mypolicy%s"
  module_id = snapcd_module.this.id

  path = "/etc/snapcd/policy-packs/aws-guardrails"
}

`)

func TestAccResourceModulePulumiLocalPolicy_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModulePulumiLocalPolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_local_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_local_policy.this", "enabled", "true"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_local_policy.this", "evaluate_on", "ApplyOnly"),
				),
			},
		},
	})
}

func TestAccResourceModulePulumiLocalPolicy_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModulePulumiLocalPolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_local_policy.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_pulumi_local_policy" "this" {
  name        = "mypolicy%s"
  module_id   = snapcd_module.this.id
  enabled     = false

  path = "/etc/snapcd/policy-packs/aws-guardrails"
}

`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_local_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_local_policy.this", "enabled", "false"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_local_policy.this", "evaluate_on", "ApplyOnly"),
				),
			},
		},
	})
}

func TestAccResourceModulePulumiLocalPolicy_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModulePulumiLocalPolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_local_policy.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_pulumi_local_policy.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
