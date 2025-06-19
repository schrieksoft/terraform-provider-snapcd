// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceBackendConfigCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_backend_config" "this" { 
  namespace_id  = snapcd_namespace.this.id
  name    		= "somevalue%s"
  value      	= "foo"
}
  
`)

func TestAccResourceNamespaceBackendConfig_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceBackendConfigCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_backend_config.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceBackendConfig_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceBackendConfigCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_backend_config.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_backend_config.this", "value", "foo"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + `
resource "snapcd_namespace_backend_config" "this" { 
  namespace_id 	= snapcd_namespace.this.id
  name  		= "somevalue%s"
  value  		= "bar"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_backend_config.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_backend_config.this", "value", "bar"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceBackendConfig_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceBackendConfigCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_backend_config.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_backend_config.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
