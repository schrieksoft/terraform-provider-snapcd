// SPDX-License-Identifier: MPL-2.0

package policies

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespacePulumiRemotePolicyCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_pulumi_remote_policy" "this" {
  name      = "mypolicy%s"
  namespace_id = snapcd_namespace.this.id

  repo_url = "https://github.com/myorg/policy-repo.git"
  revision = "v1.0.0"
  path     = "packs/aws-guardrails"
}

`)

func TestAccResourceNamespacePulumiRemotePolicy_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespacePulumiRemotePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_remote_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_remote_policy.this", "enabled", "true"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_remote_policy.this", "evaluate_on", "ApplyAndDestroy"),
				),
			},
		},
	})
}

func TestAccResourceNamespacePulumiRemotePolicy_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespacePulumiRemotePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_remote_policy.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_pulumi_remote_policy" "this" {
  name        = "mypolicy%s"
  namespace_id   = snapcd_namespace.this.id
  enabled     = false
  evaluate_on = "ApplyOnly"

  repo_url = "https://github.com/myorg/policy-repo.git"
  revision = "v1.0.0"
  path     = "packs/aws-guardrails"
}

`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_remote_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_remote_policy.this", "enabled", "false"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_remote_policy.this", "evaluate_on", "ApplyOnly"),
				),
			},
		},
	})
}

func TestAccResourceNamespacePulumiRemotePolicy_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespacePulumiRemotePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_remote_policy.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_pulumi_remote_policy.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
