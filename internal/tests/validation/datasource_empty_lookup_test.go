// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"fmt"
	"regexp"
	"testing"

	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Validation runs before the provider is configured, so these need no server. A request would be a
// test failure in itself: the point of the validator is that nothing is sent.
func TestEmptyLookupValueIsRejected(t *testing.T) {
	t.Parallel()

	cases := []struct {
		dataSource string
		attribute  string
		config     string
	}{
		{"snapcd_runner", "name", `name = ""`},
		{"snapcd_agent", "name", `name = ""`},
		{"snapcd_stack", "name", `name = ""`},
		{"snapcd_state_store", "name", `name = ""`},
		{"snapcd_integration", "name", `name = ""`},
		{"snapcd_group", "name", `name = ""`},
		{"snapcd_user", "user_name", `user_name = ""`},
		{"snapcd_service_principal", "client_id", `client_id = ""`},
		{"snapcd_namespace", "name", `stack_id = "a3bb189e-8bf9-3888-9912-ace4e6543002"
	name     = ""`},
		{"snapcd_namespace", "stack_id", `stack_id = ""
	name     = "some-namespace"`},
		{"snapcd_module", "name", `namespace_id = "a3bb189e-8bf9-3888-9912-ace4e6543002"
	name         = ""`},
		{"snapcd_module_extra_file", "file_name", `module_id = "a3bb189e-8bf9-3888-9912-ace4e6543002"
	file_name = ""`},
		{"snapcd_module_additional_trigger_path", "path", `module_id = "a3bb189e-8bf9-3888-9912-ace4e6543002"
	path      = ""`},
		{"snapcd_source_refresher_preselection", "source_url", `source_url = ""`},
		{"snapcd_module_terraform_inline_policy", "name", `module_id = "a3bb189e-8bf9-3888-9912-ace4e6543002"
	name      = ""`},
		{"snapcd_module_input_from_literal", "name", `module_id  = "a3bb189e-8bf9-3888-9912-ace4e6543002"
	input_kind = "Param"
	name       = ""`},
		{"snapcd_namespace_input_from_literal", "name", `namespace_id = "a3bb189e-8bf9-3888-9912-ace4e6543002"
	input_kind   = "Param"
	name         = ""`},
	}

	for _, test := range cases {
		t.Run(test.dataSource+"/"+test.attribute, func(t *testing.T) {
			t.Parallel()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerconfig.ProviderConfig() + fmt.Sprintf(`
data %q "this" {
	%s
}`, test.dataSource, test.config),
						// The attribute name is what the original 400 failed to mention.
						ExpectError: regexp.MustCompile(`Missing ` + test.attribute),
					},
				},
			})
		})
	}
}

// A whitespace-only value reaches the path the same way an empty one does.
func TestWhitespaceLookupValueIsRejected(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + `
data "snapcd_runner" "this" {
	name = "   "
}`,
				ExpectError: regexp.MustCompile(`Missing name`),
			},
		},
	})
}
