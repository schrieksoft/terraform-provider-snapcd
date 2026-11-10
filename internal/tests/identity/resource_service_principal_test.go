// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"fmt"
	"regexp"
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceServicePrincipal_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + ServicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "client_id", providerconfig.AppendRandomString("sp_create_%s")),
					resource.TestCheckNoResourceAttr("snapcd_service_principal.this", "client_secret"),
				),
			},
		},
	})
}

// The resource exposes no client_secret attribute, so a configuration that sets one must fail.
func TestAccResourceServicePrincipal_RejectsClientSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerconfig.ProviderConfig() + ServicePrincipalWithSecretConfig,
				ExpectError: regexp.MustCompile(`(?i)client_secret`),
			},
		},
	})
}

func TestAccResourceServicePrincipal_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + ServicePrincipalUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "is_disabled", "true"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "this" {
  client_id   = "sp_update_%s"
  is_disabled = false
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "is_disabled", "false"),
				),
			},
		},
	})
}

// client_id is RequiresReplace, so changing it must replace the Service Principal.
func TestAccResourceServicePrincipal_ReplaceOnClientIdChange(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + ServicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "this" {
  client_id = "sp_create_renamed_%s"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "client_id", providerconfig.AppendRandomString("sp_create_renamed_%s")),
				),
			},
		},
	})
}

func TestAccResourceServicePrincipal_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + ServicePrincipalImportConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
				),
			},
			{
				// Read uses the metadata endpoint, which returns neither is_disabled nor scopes,
				// so an imported state carries only the identity fields.
				ResourceName: "snapcd_service_principal.this",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 state, got %d", len(states))
					}

					want := providerconfig.AppendRandomString("sp_import_%s")
					if got := states[0].Attributes["client_id"]; got != want {
						return fmt.Errorf("client_id: want %q, got %q", want, got)
					}

					if _, ok := states[0].Attributes["client_secret"]; ok {
						return fmt.Errorf("client_secret must not appear in state")
					}

					return nil
				},
			},
		},
	})
}
