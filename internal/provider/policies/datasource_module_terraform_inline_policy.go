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

var _ datasource.DataSource = (*moduleTerraformInlinePolicyDataSource)(nil)

func ModuleTerraformInlinePolicyDataSource() datasource.DataSource {
	return &moduleTerraformInlinePolicyDataSource{}
}

type moduleTerraformInlinePolicyDataSource struct {
	client *snapcd.Client
}

func (r *moduleTerraformInlinePolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleTerraformInlinePolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_terraform_inline_policy"
}

func (d *moduleTerraformInlinePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Policies --- Use this data source to access information about an existing Module Terraform Inline Policy in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["ModuleTerraformInlinePolicy"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformInlinePolicyReadDto_Id,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleTerraformInlinePolicyReadDto_ModuleId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleTerraformInlinePolicyReadDto_Name,
			},
			"policy_content": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformInlinePolicyReadDto_PolicyContent,
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformInlinePolicyReadDto_Enabled,
			},
			"evaluate_on": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleTerraformInlinePolicyReadDto_EvaluateOn,
			},
		},
	}
}

func (d *moduleTerraformInlinePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleTerraformInlinePolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleTerraformInlinePolicyEndpoint, data.ModuleId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleTerraformInlinePolicyDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleTerraformInlinePolicyDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
