// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"strings"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecretScopedToStack_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + SecretScopedToStackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_scoped_to_stack.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceSecretScopedToStack_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + SecretScopedToStackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_scoped_to_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_secret_scoped_to_stack.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + strings.ReplaceAll(SecretScopedToStackCreateConfig, "somevalue", "someNEWvalue"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_scoped_to_stack.this", "id"),
					resource.TestCheckResourceAttr("snapcd_secret_scoped_to_stack.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceSecretScopedToStack_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + SecretScopedToStackCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_secret_scoped_to_stack.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_secret_scoped_to_stack.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
