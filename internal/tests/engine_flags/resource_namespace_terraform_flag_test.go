// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceTerraformFlagCreateConfig = `
resource "snapcd_namespace_terraform_flag" "this" {
  namespace_id 	= snapcd_namespace.this.id
  task  		= "Init"
  flag  		= "Upgrade"
}

`

func TestAccResourceNamespaceTerraformFlag_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceTerraformFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_flag.this", "task", "Init"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_flag.this", "flag", "Upgrade"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceTerraformFlag_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceTerraformFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_flag.this", "task", "Init"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_flag.this", "flag", "Upgrade"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + `
resource "snapcd_namespace_terraform_flag" "this" {
  namespace_id 	= snapcd_namespace.this.id
  task  		= "Init"
  flag  		= "Upgrade"
  value   	= "true"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_flag.this", "value", "true"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceTerraformFlag_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceTerraformFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_flag.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_terraform_flag.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
