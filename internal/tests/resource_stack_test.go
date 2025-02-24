// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var stackCreateConfig = appendRandomString(`
resource "snapcd_stack" "this" { 
  name  = "somevalue%s"
}`)

func TestAccResourceStack_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + stackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceStack_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + stackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_stack.this", "name", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + appendRandomString(`
resource "snapcd_stack" "this" { 
  name = "someNEWvalue%s"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_stack.this", "name", appendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceStack_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + stackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_stack.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
