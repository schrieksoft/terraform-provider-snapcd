// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceModuleEnvVarFromSecretScopedToNamespace(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfigDelta + ModuleEnvVarFromSecretScopedToNamespaceCreateConfig + `
data "snapcd_module_env_var_from_secret_scoped_to_namespace" "this" {
	name 		= snapcd_module_env_var_from_secret_scoped_to_namespace.this.name
	module_id 	                        = snapcd_module_env_var_from_secret_scoped_to_namespace.this.module_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_module_env_var_from_secret_scoped_to_namespace.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_module_env_var_from_secret_scoped_to_namespace.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
