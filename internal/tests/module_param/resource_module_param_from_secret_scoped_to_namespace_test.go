// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ModuleParamFromSecretScopedToNamespaceCreateConfig = secret.SimpleSecretScopedToNamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_param_from_secret_scoped_to_namespace" "this" { 
  module_id = snapcd_module.this.id
  name  	= "somevalue%s"
  secret_scoped_to_namespace_id = snapcd_simple_secret_scoped_to_namespace.this.id
}
  
`)

func TestAccResourceModuleParamFromSecretScopedToNamespace_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfigDelta + ModuleParamFromSecretScopedToNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret_scoped_to_namespace.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromSecretScopedToNamespace_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfigDelta + ModuleParamFromSecretScopedToNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret_scoped_to_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_secret_scoped_to_namespace.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfigDelta + secret.SimpleSecretScopedToNamespaceCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_module_param_from_secret_scoped_to_namespace" "this" { 
  module_id = snapcd_module.this.id
  name  = "somevalue%s"
  secret_scoped_to_namespace_id = snapcd_simple_secret_scoped_to_namespace.this.id
  type = "NotString"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret_scoped_to_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_param_from_secret_scoped_to_namespace.this", "type", "NotString"),
				),
			},
		},
	})
}

func TestAccResourceModuleParamFromSecretScopedToNamespace_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.ModuleCreateConfigDelta + ModuleParamFromSecretScopedToNamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_param_from_secret_scoped_to_namespace.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_param_from_secret_scoped_to_namespace.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
