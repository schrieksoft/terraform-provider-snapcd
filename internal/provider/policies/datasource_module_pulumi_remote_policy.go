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

var _ datasource.DataSource = (*modulePulumiRemotePolicyDataSource)(nil)

func ModulePulumiRemotePolicyDataSource() datasource.DataSource {
	return &modulePulumiRemotePolicyDataSource{}
}

type modulePulumiRemotePolicyDataSource struct {
	client *snapcd.Client
}

func (r *modulePulumiRemotePolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *modulePulumiRemotePolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_pulumi_remote_policy"
}

func (d *modulePulumiRemotePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Policies --- Use this data source to access information about an existing Module Pulumi Remote Policy in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["ModulePulumiRemotePolicy"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModulePulumiRemotePolicyReadDto_Id,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModulePulumiRemotePolicyReadDto_ModuleId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModulePulumiRemotePolicyReadDto_Name,
			},
			"repo_url": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModulePulumiRemotePolicyReadDto_RepoUrl,
			},
			"revision": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModulePulumiRemotePolicyReadDto_Revision,
			},
			"path": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModulePulumiRemotePolicyReadDto_Path,
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModulePulumiRemotePolicyReadDto_Enabled,
			},
			"evaluate_on": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModulePulumiRemotePolicyReadDto_EvaluateOn,
			},
		},
	}
}

func (d *modulePulumiRemotePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data modulePulumiRemotePolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", modulePulumiRemotePolicyEndpoint, data.ModuleId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(modulePulumiRemotePolicyDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(modulePulumiRemotePolicyDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
