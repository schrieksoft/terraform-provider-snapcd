package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ExternallyDeleted drops a resource from state after the server reported it gone, and says so.
//
// Removing it silently leaves the operator with a plan that recreates something they did not
// touch and no indication why.
func ExternallyDeleted(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	resourceType string,
	id string,
) {
	tflog.Info(ctx, "resource no longer exists on the server, removing from state", map[string]any{
		"resource_type": resourceType,
		"id":            id,
	})

	diagnostics.AddWarning(
		fmt.Sprintf("%s was deleted outside Terraform", resourceType),
		fmt.Sprintf(
			"%s with ID %s no longer exists on the server, so it has been removed from state "+
				"and will be recreated.",
			resourceType, id),
	)
}
