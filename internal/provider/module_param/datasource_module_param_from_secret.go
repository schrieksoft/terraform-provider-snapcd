package module_param

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*moduleParamFromSecretDataSource)(nil)

func ModuleParamFromSecretDataSource() datasource.DataSource {
	return &moduleParamFromSecretDataSource{}
}

type moduleParamFromSecretDataSource struct {
	client *snapcd.Client
}

func (r *moduleParamFromSecretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleParamFromSecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_param_from_secret"
}

func (d *moduleParamFromSecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Module Inputs --- Use this data source to access information about an existing Module Param (From Secret) in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedId + "Module Param (From Secret).",
			},
			"name": schema.StringAttribute{
				Required: true,
				Description: DescSharedName1 + "Module Param (From Secret). " + DescSharedName2,
			},
			"module_id": schema.StringAttribute{
				Required: true,
				Description: DescSharedModuleId1 + "Module Param (From Secret)" + DescSharedModuleId2,
			},
			"type": schema.StringAttribute{
				Computed: true,
				Description: DescSharedSecretType,
			},
			"secret_name": schema.StringAttribute{
				Computed: true,
				Description: DescSharedSecretName,
			},
			"secret_scope": schema.StringAttribute{
				Computed: true,
				Description: DescSharedSecretScope,
			},
		},
	}
}

func (d *moduleParamFromSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleParamFromSecretModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleParamFromSecretEndpoint, data.ModuleId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
