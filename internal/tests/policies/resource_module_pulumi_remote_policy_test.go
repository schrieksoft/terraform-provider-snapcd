// SPDX-License-Identifier: MPL-2.0

package policies

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModulePulumiRemotePolicyCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_pulumi_remote_policy" "this" {
  name      = "mypolicy%s"
  module_id = snapcd_module.this.id

  repo_url = "https://github.com/myorg/policy-repo.git"
  revision = "v1.0.0"
  path     = "packs/aws-guardrails"
}

`)

func TestAccResourceModulePulumiRemotePolicy_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModulePulumiRemotePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_remote_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_remote_policy.this", "enabled", "true"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_remote_policy.this", "evaluate_on", "ApplyOnly"),
				),
			},
		},
	})
}

func TestAccResourceModulePulumiRemotePolicy_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModulePulumiRemotePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_remote_policy.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_pulumi_remote_policy" "this" {
  name        = "mypolicy%s"
  module_id   = snapcd_module.this.id
  enabled     = false

  repo_url = "https://github.com/myorg/policy-repo.git"
  revision = "v1.0.0"
  path     = "packs/aws-guardrails"
}

`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_remote_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_remote_policy.this", "enabled", "false"),
					resource.TestCheckResourceAttr("snapcd_module_pulumi_remote_policy.this", "evaluate_on", "ApplyOnly"),
				),
			},
		},
	})
}

func TestAccResourceModulePulumiRemotePolicy_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModulePulumiRemotePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_pulumi_remote_policy.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_pulumi_remote_policy.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
