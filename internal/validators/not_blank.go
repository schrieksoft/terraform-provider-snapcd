package validators

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// notBlankValidator rejects a string that is empty or only whitespace.
type notBlankValidator struct{}

func (v notBlankValidator) Description(_ context.Context) string {
	return "value must not be empty or only whitespace"
}

func (v notBlankValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v notBlankValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// An unknown value belongs to a resource that has not been applied yet, so there is nothing to judge.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if strings.TrimSpace(req.ConfigValue.ValueString()) != "" {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Missing "+req.Path.String(),
		"The value of "+req.Path.String()+" is empty. This usually means it was taken from a variable or an "+
			"attribute that resolved to an empty string. Give it a value that identifies the object to look up.",
	)
}

// NotBlank returns a validator that rejects an empty or whitespace-only string.
//
// Lookup values are interpolated into the request path, where an empty value silently changes which endpoint
// the request reaches.
func NotBlank() validator.String {
	return notBlankValidator{}
}
