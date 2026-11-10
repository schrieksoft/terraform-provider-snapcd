// SPDX-License-Identifier: MPL-2.0

package identity

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const federatedCredential = "snapcd_service_principal_federated_credential.this"

func TestAccResourceServicePrincipalFederatedCredential_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + FederatedCredentialCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(federatedCredential, "id"),
					resource.TestCheckResourceAttrPair(
						federatedCredential, "service_principal_id",
						"snapcd_service_principal.fedcred", "id"),
					resource.TestCheckResourceAttr(federatedCredential, "name", "fedcred_create"),
					resource.TestCheckResourceAttr(federatedCredential, "issuer", "https://oidc.example.com/cluster"),
					resource.TestCheckResourceAttr(federatedCredential, "subject", "system:serviceaccount:snapcd:runner"),
					resource.TestCheckResourceAttr(federatedCredential, "audience", "snapcd"),
				),
			},
		},
	})
}

func TestAccResourceServicePrincipalFederatedCredential_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + FederatedCredentialCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(federatedCredential, "id"),
					resource.TestCheckResourceAttr(federatedCredential, "subject", "system:serviceaccount:snapcd:runner"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + FederatedCredentialUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(federatedCredential, "id"),
					resource.TestCheckResourceAttr(federatedCredential, "issuer", "https://oidc.example.com/other-cluster"),
					resource.TestCheckResourceAttr(federatedCredential, "subject", "system:serviceaccount:snapcd:agent"),
					resource.TestCheckResourceAttr(federatedCredential, "is_disabled", "true"),
				),
			},
		},
	})
}

// service_principal_id is RequiresReplace, so pointing the credential at another Service
// Principal must replace it.
func TestAccResourceServicePrincipalFederatedCredential_ReplaceOnServicePrincipalChange(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + FederatedCredentialCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(federatedCredential, "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + FederatedCredentialReplaceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						federatedCredential, "service_principal_id",
						"snapcd_service_principal.other", "id"),
				),
			},
		},
	})
}

func TestAccResourceServicePrincipalFederatedCredential_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + FederatedCredentialCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(federatedCredential, "id"),
				),
			},
			{
				ResourceName:      federatedCredential,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
