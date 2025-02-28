// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var namespaceCreateConfig = appendRandomString(`
data "snapcd_stack" "default" { 
  name  = "default"
}

resource "snapcd_namespace" "this" { 
  name                        	 = "somevalue%s"
  stack_id			     		 = data.snapcd_stack.default.id
  default_init_before_hook       = "foo"
}
`)

var namespaceUpdateConfig = appendRandomString(`


data "snapcd_stack" "default" { 
  name  = "default"
}

resource "snapcd_namespace" "this" { 
  name                        	 = "somevalue%s"
  stack_id			     		 = data.snapcd_stack.default.id
  default_init_before_hook       = "bar"
}

`)

func TestAccResourceNamespace_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespace_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace.this", "default_init_before_hook", "foo"),
				),
			},
			{
				Config: providerConfig + namespaceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace.this", "default_init_before_hook", "bar"),
				),
			},
		},
	})
}

func TestAccResourceNamespace_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + namespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
