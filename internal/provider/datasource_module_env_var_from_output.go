package provider

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*moduleEnvVarFromOutputDataSource)(nil)

func ModuleEnvVarFromOutputDataSource() datasource.DataSource {
	return &moduleEnvVarFromOutputDataSource{}
}

type moduleEnvVarFromOutputDataSource struct {
	client *snapcd.Client
}

func (r *moduleEnvVarFromOutputDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleEnvVarFromOutputDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_env_var_from_output"
}

func (d *moduleEnvVarFromOutputDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"output_name": schema.StringAttribute{
				Computed: true,
			},
			"module_name": schema.StringAttribute{
				Computed: true,
			},
			"namespace_name": schema.StringAttribute{
				Computed: true,
			},
			"module_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (d *moduleEnvVarFromOutputDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleEnvVarFromOutputModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleEnvVarFromOutputEndpoint, data.ModuleId.ValueString(), data.Name.ValueString()))
	err := httpError.Error

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromOutputDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromOutputDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
