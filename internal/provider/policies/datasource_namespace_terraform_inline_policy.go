package policies

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*namespaceTerraformInlinePolicyDataSource)(nil)

func NamespaceTerraformInlinePolicyDataSource() datasource.DataSource {
	return &namespaceTerraformInlinePolicyDataSource{}
}

type namespaceTerraformInlinePolicyDataSource struct {
	client *snapcd.Client
}

func (r *namespaceTerraformInlinePolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*snapcd.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *snapcd.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (d *namespaceTerraformInlinePolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_terraform_inline_policy"
}

func (d *namespaceTerraformInlinePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Policies --- Use this data source to access information about an existing Namespace Terraform Inline Policy in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourceMetadataPermissions["NamespaceTerraformInlinePolicy"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespaceTerraformInlinePolicyMetadataReadDto_Id,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceTerraformInlinePolicyMetadataReadDto_NamespaceId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceTerraformInlinePolicyMetadataReadDto_Name,
			},
		},
	}
}

func (d *namespaceTerraformInlinePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceTerraformInlinePolicyMetadataModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/Metadata/%s/%s", namespaceTerraformInlinePolicyEndpoint, data.NamespaceId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformInlinePolicyDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformInlinePolicyDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Metadata only: a data source exists to resolve an ID, and narrowing it means a
// principal who may only discover the resource can still use it.
type namespaceTerraformInlinePolicyMetadataModel struct {
	Id          types.String `tfsdk:"id"`
	NamespaceId types.String `tfsdk:"namespace_id"`
	Name        types.String `tfsdk:"name"`
}
