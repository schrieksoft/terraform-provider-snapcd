// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNamespace(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `

data "snapcd_stack" "default" {
  name = "default"
}  

data "snapcd_namespace" "default" {
	name 	 = "default"
	stack_id = data.snapcd_stack.default.id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_namespace.default", "id"),
					resource.TestCheckResourceAttr("data.snapcd_namespace.default", "name", "default"),
				),
			},
		},
	})
}
