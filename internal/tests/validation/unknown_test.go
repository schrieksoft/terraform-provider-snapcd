// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"os"
	"regexp"
	"testing"

	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The reported case: a name taken from a variable that resolved to empty.
//
// Terraform configures the provider before it validates a data source whose value comes from a
// variable, so this needs a reachable server; with none, the client error arrives first and the
// validator is never reached. The literal-value cases above cover the validator itself.
func TestEmptyFromVariableIsRejected(t *testing.T) {
	t.Parallel()

	if os.Getenv("SNAPCD_URL") == "" {
		t.Skip("needs a reachable server: set SNAPCD_URL")
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + `
variable "runner_name" {
	type    = string
	default = ""
}

data "snapcd_runner" "app" {
	name = var.runner_name
}`,
				ExpectError: regexp.MustCompile(`Missing name`),
			},
		},
	})
}

// A value that is unknown at plan time belongs to a resource not yet created; validation must not
// reject it or ordinary configurations break.
func TestUnknownValuePassesValidation(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + `
resource "terraform_data" "pending" {
	input = "some-runner"
}

data "snapcd_runner" "app" {
	name = terraform_data.pending.output
}`,
				// Validation must pass, so this gets as far as the request and fails there instead.
				ExpectError: regexp.MustCompile(`(?s)snapcd_runner error|connection refused|Post "|health`),
			},
		},
	})
}
