// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The token must come from an issuer the Server can reach over https, so these run only where
// one has been minted: SNAPCD_TEST_FEDERATED_TOKEN_FILE names it.
func federatedTokenFile(t *testing.T) string {
	path := os.Getenv("SNAPCD_TEST_FEDERATED_TOKEN_FILE")
	if path == "" {
		t.Skip("SNAPCD_TEST_FEDERATED_TOKEN_FILE is not set")
	}
	return path
}

func TestAccProviderFederatedToken(t *testing.T) {
	tokenFile := federatedTokenFile(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.FederatedProviderConfig(tokenFile) + ServicePrincipalDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_service_principal.this", "id"),
				),
			},
		},
	})
}

func TestAccProviderFederatedTokenTampered(t *testing.T) {
	content, err := os.ReadFile(federatedTokenFile(t))
	if err != nil {
		t.Fatal(err)
	}

	// Same header and claims, with the signature replaced, so only the signature check can refuse it.
	parts := strings.Split(strings.TrimSpace(string(content)), ".")
	if len(parts) != 3 {
		t.Fatalf("token file does not hold a JWT")
	}
	parts[2] = strings.Repeat("A", len(parts[2]))
	tampered := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(tampered, []byte(strings.Join(parts, ".")), 0o600); err != nil {
		t.Fatal(err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerconfig.FederatedProviderConfig(tampered) + ServicePrincipalDataSourceConfig,
				ExpectError: regexp.MustCompile(`authentication failed with\s+status 401`),
			},
		},
	})
}
