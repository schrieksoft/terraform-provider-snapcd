// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceAwsSecretsManagerSecretScopedToModule(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + AwsSecretsManagerSecretScopedToModuleCreateConfig + `

data "snapcd_aws_secrets_manager_secret_scoped_to_module" "this" {
	name = snapcd_aws_secrets_manager_secret_scoped_to_module.this.name
    module_id = snapcd_module.this.id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_aws_secrets_manager_secret_scoped_to_module.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_aws_secrets_manager_secret_scoped_to_module.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
