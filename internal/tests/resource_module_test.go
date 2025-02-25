// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var moduleCreateConfig = namespaceCreateConfig + appendRandomString(`

data "snapcd_runner_pool" "default" {
  name = "default"
}

resource "snapcd_module" "this" { 
  name                         	 = "somevalue%s"
  namespace_id                	 = snapcd_namespace.this.id
  runner_pool_id                 = data.snapcd_runner_pool.default.id
// target_repo_revision         	 = "main"
//   target_repo_url              	 = "git@github.com:karlschriek/tf-samples.git"
  target_module_relative_path  	 = "modules/module1"
  provider_cache_enabled         = true
  module_cache_enabled         	 = true
  depends_on_modules		 	 = []
  select_on           			 = "PoolId"
  select_strategy     			 = "FirstOf"
  init_before_hook				 = "fooBeforeHook"
  //selected_consumer_id 		 = "cli"
}
`)

func TestAccResourceModule_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceModule_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module.this", "init_before_hook", "fooBeforeHook"),
				),
			},
			{
				Config: providerConfig + strings.Replace(moduleCreateConfig, "fooBeforeHook", "barBeforeHook", -1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module.this", "init_before_hook", "barBeforeHook"),
				),
			},
		},
	})
}

func TestAccResourceModule_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + moduleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
