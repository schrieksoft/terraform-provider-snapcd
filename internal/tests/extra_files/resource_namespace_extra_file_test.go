// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceExtraFileCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_namespace_extra_file" "this" { 
  namespace_id = snapcd_namespace.this.id
  file_name    = "somevalue%s"
  contents     = "foo"
}
  
`)

func TestAccResourceNamespaceExtraFile_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceExtraFileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_extra_file.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceExtraFile_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceExtraFileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_extra_file.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_extra_file.this", "contents", "foo"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + `
resource "snapcd_namespace_extra_file" "this" { 
  namespace_id = snapcd_namespace.this.id
  file_name  = "somevalue%s"
  contents  = "bar"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_extra_file.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_extra_file.this", "contents", "bar"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceExtraFile_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceExtraFileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_extra_file.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_extra_file.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
