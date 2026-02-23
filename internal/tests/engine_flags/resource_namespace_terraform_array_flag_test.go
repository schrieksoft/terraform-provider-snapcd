// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceTerraformArrayFlagCreateConfig = `
resource "snapcd_namespace_terraform_array_flag" "this" {
  namespace_id 	= snapcd_namespace.this.id
  task  		= "Plan"
  flag  		= "Target"
  value  		= "aws_vpc.main"
}

`

func TestAccResourceNamespaceTerraformArrayFlag_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceTerraformArrayFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_array_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_array_flag.this", "task", "Plan"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_array_flag.this", "flag", "Target"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceTerraformArrayFlag_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceTerraformArrayFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_array_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_array_flag.this", "value", "aws_vpc.main"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + `
resource "snapcd_namespace_terraform_array_flag" "this" {
  namespace_id 	= snapcd_namespace.this.id
  task  		= "Plan"
  flag  		= "Target"
  value   	= "aws_subnet.main"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_array_flag.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_terraform_array_flag.this", "value", "aws_subnet.main"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceTerraformArrayFlag_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceTerraformArrayFlagCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_terraform_array_flag.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_terraform_array_flag.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
