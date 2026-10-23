// SPDX-License-Identifier: MPL-2.0

package validators_test

import (
	"context"
	"testing"

	"terraform-provider-snapcd/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNotBlank(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value     types.String
		wantError bool
	}{
		"a name":         {types.StringValue("some-runner"), false},
		"inner spaces":   {types.StringValue("some runner"), false},
		"padded":         {types.StringValue("  some-runner  "), false},
		"empty":          {types.StringValue(""), true},
		"only spaces":    {types.StringValue("   "), true},
		"tab":            {types.StringValue("\t"), true},
		"newline":        {types.StringValue("\n"), true},
		"null passes":    {types.StringNull(), false},
		"unknown passes": {types.StringUnknown(), false},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &validator.StringResponse{}
			validators.NotBlank().ValidateString(
				context.Background(),
				validator.StringRequest{
					Path:        path.Root("name"),
					ConfigValue: test.value,
				},
				resp,
			)

			if got := resp.Diagnostics.HasError(); got != test.wantError {
				t.Fatalf("HasError() = %v, want %v (diagnostics: %v)", got, test.wantError, resp.Diagnostics)
			}

			if !test.wantError {
				return
			}

			// The whole point of the validator is that the message names the attribute.
			if summary := resp.Diagnostics.Errors()[0].Summary(); summary != "Missing name" {
				t.Errorf("summary = %q, want it to name the attribute", summary)
			}
		})
	}
}
