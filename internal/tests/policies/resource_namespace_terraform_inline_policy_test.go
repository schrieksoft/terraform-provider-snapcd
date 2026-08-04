// SPDX-License-Identifier: MPL-2.0

package policies

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceTerraformInlinePolicyCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_terraform_inline_policy" "this" {
  name      = "mypolicy%s"
  namespace_id = snapcd_namespace.this.id

  policy_content = <<-EOF
    package snapcd

    import rego.v1

    deny contains msg if {
      false
      msg := "never fires"
    }
  EOF
}

`)

func TestAccResourceNamespaceTerraformInlinePolicy_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceTerraformInlinePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_inline_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_inline_policy.this", "enabled", "true"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_inline_policy.this", "evaluate_on", "ApplyAndDestroy"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceTerraformInlinePolicy_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceTerraformInlinePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_inline_policy.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_terraform_inline_policy" "this" {
  name        = "mypolicy%s"
  namespace_id   = snapcd_namespace.this.id
  enabled     = false
  evaluate_on = "ApplyOnly"

  policy_content = <<-EOF
    package snapcd

    import rego.v1

    deny contains msg if {
      false
      msg := "never fires"
    }
  EOF
}

`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_inline_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_inline_policy.this", "enabled", "false"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_inline_policy.this", "evaluate_on", "ApplyOnly"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceTerraformInlinePolicy_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceTerraformInlinePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_inline_policy.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_terraform_inline_policy.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
