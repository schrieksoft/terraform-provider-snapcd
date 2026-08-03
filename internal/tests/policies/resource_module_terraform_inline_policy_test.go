// SPDX-License-Identifier: MPL-2.0

package policies

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleTerraformInlinePolicyCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_module_terraform_inline_policy" "this" {
  name      = "mypolicy%s"
  module_id = snapcd_module.this.id

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

func TestAccResourceModuleTerraformInlinePolicy_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleTerraformInlinePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_terraform_inline_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_inline_policy.this", "enabled", "true"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_inline_policy.this", "evaluate_on", "ApplyAndDestroy"),
				),
			},
		},
	})
}

func TestAccResourceModuleTerraformInlinePolicy_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleTerraformInlinePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_terraform_inline_policy.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_terraform_inline_policy" "this" {
  name        = "mypolicy%s"
  module_id   = snapcd_module.this.id
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
					resource.TestCheckResourceAttrSet("snapcd_module_terraform_inline_policy.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_inline_policy.this", "enabled", "false"),
					resource.TestCheckResourceAttr("snapcd_module_terraform_inline_policy.this", "evaluate_on", "ApplyOnly"),
				),
			},
		},
	})
}

func TestAccResourceModuleTerraformInlinePolicy_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.ModuleCreateConfig + ModuleTerraformInlinePolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_terraform_inline_policy.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_terraform_inline_policy.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
