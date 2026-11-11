// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"os"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Needs a snapcd CLI logged in to the test Server, so it runs only when SNAPCD_TEST_CLI is set.
// The provider block names no Server, organization or credentials: all of them come from the CLI.
func TestAccProviderCliLogin(t *testing.T) {
	if os.Getenv("SNAPCD_TEST_CLI") == "" {
		t.Skip("SNAPCD_TEST_CLI is not set")
	}
	for _, name := range []string{"SNAPCD_URL", "SNAPCD_ORGANIZATION_ID", "SNAPCD_ACCESS_TOKEN", "SNAPCD_CLIENT_ID", "SNAPCD_CLIENT_SECRET", "SNAPCD_FEDERATED_TOKEN_FILE"} {
		t.Setenv(name, "")
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "snapcd" {
	health_check_interval_seconds = 5
	health_check_timeout_seconds  = 20
}
` + ServicePrincipalDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snapcd_service_principal.this", "id"),
				),
			},
		},
	})
}
