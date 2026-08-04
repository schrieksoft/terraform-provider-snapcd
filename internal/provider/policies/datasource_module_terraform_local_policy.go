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

var _ datasource.DataSource = (*moduleTerraformLocalPolicyDataSource)(nil)

func ModuleTerraformLocalPolicyDataSource() datasource.DataSource {
	return &moduleTerraformLocalPolicyDataSource{}
}

type moduleTerraformLocalPolicyDataSource struct {
	client *snapcd.Client
}

func (r *moduleTerraformLocalPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleTerraformLocalPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_terraform_local_policy"
}

func (d *moduleTerraformLocalPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Policies --- Use this data source to access information about an existing Module Terraform Local Policy in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["ModuleTerraformLocalPolicy"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformLocalPolicyReadDto_Id,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleTerraformLocalPolicyReadDto_ModuleId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleTerraformLocalPolicyReadDto_Name,
			},
			"path": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformLocalPolicyReadDto_Path,
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformLocalPolicyReadDto_Enabled,
			},
			"evaluate_on": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformLocalPolicyReadDto_EvaluateOn,
			},
		},
	}
}

func (d *moduleTerraformLocalPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleTerraformLocalPolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleTerraformLocalPolicyEndpoint, data.ModuleId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleTerraformLocalPolicyDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleTerraformLocalPolicyDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
