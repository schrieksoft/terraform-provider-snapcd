// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceBackendConfig(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceBackendConfigCreateConfig + `
data "snapcd_namespace_backend_config" "this" {
	name 			= snapcd_namespace_backend_config.this.name
	namespace_id 	= snapcd_namespace_backend_config.this.namespace_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_backend_config.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_backend_config.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
