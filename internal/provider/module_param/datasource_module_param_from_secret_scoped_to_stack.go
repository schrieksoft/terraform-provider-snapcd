package module_param

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*moduleParamFromSecretScopedToStackDataSource)(nil)

func ModuleParamFromSecretScopedToStackDataSource() datasource.DataSource {
	return &moduleParamFromSecretScopedToStackDataSource{}
}

type moduleParamFromSecretScopedToStackDataSource struct {
	client *snapcd.Client
}

func (r *moduleParamFromSecretScopedToStackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleParamFromSecretScopedToStackDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_param_from_secret_scoped_to_stack"
}

func (d *moduleParamFromSecretScopedToStackDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Module Inputs (Parameters) --- Use this data source to access information about an existing Module Param (From Secret Scoped To Stack) in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedId + "Module Param (From Secret Scoped To Stack).",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedName1 + "Module Param (From Secret Scoped To Stack). " + DescSharedName2,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedModuleId1 + "Module Param (From Secret Scoped To Stack)" + DescSharedModuleId2,
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedSecretType,
			},
			"secret_scoped_to_stack_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the Secret Scoped To Stack to take as input.",
			},
		},
	}
}

func (d *moduleParamFromSecretScopedToStackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleParamFromSecretScopedToStackModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleParamFromSecretScopedToStackEndpoint, data.ModuleId.ValueString(), data.SecretScopedToStackId.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToStackDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToStackDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
