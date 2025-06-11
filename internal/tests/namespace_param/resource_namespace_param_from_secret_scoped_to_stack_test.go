// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/secret"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceParamFromSecretScopedToStackCreateConfig = secret.SimpleSecretScopedToStackCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_param_from_secret_scoped_to_stack" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  	= "somevalue%s"
  secret_scoped_to_stack_id = snapcd_simple_secret_scoped_to_stack.this.id
}
  
`)

func TestAccResourceNamespaceParamFromSecretScopedToStack_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceParamFromSecretScopedToStackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_secret_scoped_to_stack.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceParamFromSecretScopedToStack_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceParamFromSecretScopedToStackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_secret_scoped_to_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_param_from_secret_scoped_to_stack.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + secret.SimpleSecretScopedToStackCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace_param_from_secret_scoped_to_stack" "this" { 
  namespace_id = snapcd_namespace.this.id
  name  = "somevalue%s"
  secret_scoped_to_stack_id = snapcd_simple_secret_scoped_to_stack.this.id
  usage_mode = "UseByDefault"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_secret_scoped_to_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_param_from_secret_scoped_to_stack.this", "usage_mode", "UseByDefault"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceParamFromSecretScopedToStack_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceParamFromSecretScopedToStackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_param_from_secret_scoped_to_stack.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_param_from_secret_scoped_to_stack.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
