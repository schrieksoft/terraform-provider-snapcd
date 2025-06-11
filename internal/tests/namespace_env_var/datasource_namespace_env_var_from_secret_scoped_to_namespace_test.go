// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceEnvVarFromSecretScopedToNamespace(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceEnvVarFromSecretScopedToNamespaceCreateConfig + `
data "snapcd_namespace_env_var_from_secret_scoped_to_namespace" "this" {
	secret_scoped_to_namespace_id 			= snapcd_namespace_env_var_from_secret_scoped_to_namespace.this.secret_scoped_to_namespace_id
	namespace_id 	= snapcd_namespace_env_var_from_secret_scoped_to_namespace.this.namespace_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_env_var_from_secret_scoped_to_namespace.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_env_var_from_secret_scoped_to_namespace.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
