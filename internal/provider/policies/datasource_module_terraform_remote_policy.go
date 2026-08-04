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

var _ datasource.DataSource = (*moduleTerraformRemotePolicyDataSource)(nil)

func ModuleTerraformRemotePolicyDataSource() datasource.DataSource {
	return &moduleTerraformRemotePolicyDataSource{}
}

type moduleTerraformRemotePolicyDataSource struct {
	client *snapcd.Client
}

func (r *moduleTerraformRemotePolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleTerraformRemotePolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_terraform_remote_policy"
}

func (d *moduleTerraformRemotePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Policies --- Use this data source to access information about an existing Module Terraform Remote Policy in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["ModuleTerraformRemotePolicy"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformRemotePolicyReadDto_Id,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleTerraformRemotePolicyReadDto_ModuleId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleTerraformRemotePolicyReadDto_Name,
			},
			"repo_url": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformRemotePolicyReadDto_RepoUrl,
			},
			"revision": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformRemotePolicyReadDto_Revision,
			},
			"path": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformRemotePolicyReadDto_Path,
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformRemotePolicyReadDto_Enabled,
			},
			"evaluate_on": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformRemotePolicyReadDto_EvaluateOn,
			},
		},
	}
}

func (d *moduleTerraformRemotePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleTerraformRemotePolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleTerraformRemotePolicyEndpoint, data.ModuleId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleTerraformRemotePolicyDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleTerraformRemotePolicyDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
