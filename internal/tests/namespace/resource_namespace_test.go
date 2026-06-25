// SPDX-License-Identifier: MPL-2.0

package namespace

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceNamespace_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespace_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace.this", "default_apply_approval_threshold", "1"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + NamespaceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace.this", "default_apply_approval_threshold", "2"),
				),
			},
		},
	})
}

func TestAccResourceNamespace_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + NamespaceCreateConfig,
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
