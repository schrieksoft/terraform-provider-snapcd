// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespaceExtraFile(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceExtraFileCreateConfig + `
data "snapcd_namespace_extra_file" "this" {
	file_name 		= snapcd_namespace_extra_file.this.file_name
	namespace_id 	= snapcd_namespace_extra_file.this.namespace_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace_extra_file.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace_extra_file.this", "file_name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
		},
	})
}
