// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceStack(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + stackCreateConfig + `

data "snapcd_stack" "this" {
	name = snapcd_stack.this.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_stack.this", "id"),
					resource.TestCheckResourceAttr("data.snapcd_stack.this", "name", appendRandomString("somevalue%s")),
				),
			},
		},
	})
}
