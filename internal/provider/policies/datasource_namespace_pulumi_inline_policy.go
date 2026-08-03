package policies

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*namespacePulumiInlinePolicyDataSource)(nil)

func NamespacePulumiInlinePolicyDataSource() datasource.DataSource {
	return &namespacePulumiInlinePolicyDataSource{}
}

type namespacePulumiInlinePolicyDataSource struct {
	client *snapcd.Client
}

func (r *namespacePulumiInlinePolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespacePulumiInlinePolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_pulumi_inline_policy"
}

func (d *namespacePulumiInlinePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Policies --- Use this data source to access information about an existing Namespace Pulumi Inline Policy in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["NamespacePulumiInlinePolicy"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespacePulumiInlinePolicyReadDto_Id,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespacePulumiInlinePolicyReadDto_NamespaceId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespacePulumiInlinePolicyReadDto_Name,
			},
			"policy_content": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespacePulumiInlinePolicyReadDto_PolicyContent,
			},
			"runtime": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespacePulumiInlinePolicyReadDto_Runtime,
			},
			"additional_dependencies": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespacePulumiInlinePolicyReadDto_AdditionalDependencies,
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.NamespacePulumiInlinePolicyReadDto_Enabled,
			},
			"evaluate_on": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespacePulumiInlinePolicyReadDto_EvaluateOn,
			},
		},
	}
}

func (d *namespacePulumiInlinePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespacePulumiInlinePolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", namespacePulumiInlinePolicyEndpoint, data.NamespaceId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(namespacePulumiInlinePolicyDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespacePulumiInlinePolicyDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
