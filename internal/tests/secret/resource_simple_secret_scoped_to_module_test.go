// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"strings"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSimpleSecretScopedToModule_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + SimpleSecretScopedToModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_simple_secret_scoped_to_module.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceSimpleSecretScopedToModule_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + SimpleSecretScopedToModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_simple_secret_scoped_to_module.this", "id"),
					resource.TestCheckResourceAttr("snapcd_simple_secret_scoped_to_module.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{

				Config: providerconfig.ProviderConfig + strings.ReplaceAll(SimpleSecretScopedToModuleCreateConfig, "somevalue", "someNEWvalue"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_simple_secret_scoped_to_module.this", "id"),
					resource.TestCheckResourceAttr("snapcd_simple_secret_scoped_to_module.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceSimpleSecretScopedToModule_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + SimpleSecretScopedToModuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_simple_secret_scoped_to_module.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_simple_secret_scoped_to_module.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
