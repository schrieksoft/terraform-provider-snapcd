// SPDX-License-Identifier: MPL-2.0

package policies

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceTerraformLocalPolicy(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + NamespaceTerraformLocalPolicyCreateConfig + `
data "snapcd_namespace_terraform_local_policy" "this" {
	name      = snapcd_namespace_terraform_local_policy.this.name
	namespace_id = snapcd_namespace_terraform_local_policy.this.namespace_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_terraform_local_policy.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_terraform_local_policy.this", "name", providerconfig.AppendRandomString("mypolicy%s")),
					resource.TestCheckResourceAttr("data.snapcd_namespace_terraform_local_policy.this", "evaluate_on", "ApplyAndDestroy"),
				),
			},
		},
	})
}
