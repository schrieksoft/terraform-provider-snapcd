// SPDX-License-Identifier: MPL-2.0

package core

import (
	"regexp"
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStack_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + StackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceStack_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + StackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_stack.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_stack" "this" { 
  name = "someNEWvalue%s"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_stack.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceStack_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + StackCreateConfig,
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

func TestAccResourceStack_CreateShouldFail(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Test for any error (catch-all)
				Config:      providerconfig.ProviderConfig + PrexistingStack,
				ExpectError: regexp.MustCompile(`.*`), // Matches any error
			},
		},
	})
}
