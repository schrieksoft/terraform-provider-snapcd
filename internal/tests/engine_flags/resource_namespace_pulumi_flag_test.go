// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespacePulumiFlagCreateConfig = `
resource "snapcd_namespace_pulumi_flag" "this" {
  namespace_id  = snapcd_namespace.this.id
  task  		= "Init"
  flag  		= "Refresh"
  value  		= "true"
}

`

func TestAccResourceNamespacePulumiFlag_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespacePulumiFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_flag.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespacePulumiFlag_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespacePulumiFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_flag.this", "value", "true"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_flag.this", "task", "Init"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_flag.this", "flag", "Refresh"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + `
resource "snapcd_namespace_pulumi_flag" "this" {
  namespace_id 	= snapcd_namespace.this.id
  task  		= "Init"
  flag  		= "Refresh"
  value  		= "false"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_pulumi_flag.this", "value", "false"),
				),
			},
		},
	})
}

func TestAccResourceNamespacePulumiFlag_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespacePulumiFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_pulumi_flag.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_pulumi_flag.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
