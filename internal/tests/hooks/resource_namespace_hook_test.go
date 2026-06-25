// SPDX-License-Identifier: MPL-2.0

package hooks

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceHookCreateConfig = `
resource "snapcd_namespace_hook" "this" {
  namespace_id 	= snapcd_namespace.this.id
  task  		= "Plan"
  phase 		= "Before"
  script  		= "echo 'namespace before plan'"
}

`

func TestAccResourceNamespaceHook_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + NamespaceHookCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_hook.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_hook.this", "task", "Plan"),
					resource.TestCheckResourceAttr("snapcd_namespace_hook.this", "phase", "Before"),
					resource.TestCheckResourceAttr("snapcd_namespace_hook.this", "script", "echo 'namespace before plan'"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceHook_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + NamespaceHookCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_hook.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_hook.this", "task", "Plan"),
					resource.TestCheckResourceAttr("snapcd_namespace_hook.this", "phase", "Before"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + `
resource "snapcd_namespace_hook" "this" {
  namespace_id 	= snapcd_namespace.this.id
  task  		= "Plan"
  phase 		= "Before"
  script  		= "echo 'updated namespace script'"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_hook.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_hook.this", "script", "echo 'updated namespace script'"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceHook_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.NamespaceCreateConfig + NamespaceHookCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_hook.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_hook.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
