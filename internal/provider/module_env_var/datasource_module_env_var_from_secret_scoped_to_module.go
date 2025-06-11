package module_env_var

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*moduleEnvVarFromSecretScopedToModuleDataSource)(nil)

func ModuleEnvVarFromSecretScopedToModuleDataSource() datasource.DataSource {
	return &moduleEnvVarFromSecretScopedToModuleDataSource{}
}

type moduleEnvVarFromSecretScopedToModuleDataSource struct {
	client *snapcd.Client
}

func (r *moduleEnvVarFromSecretScopedToModuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleEnvVarFromSecretScopedToModuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_env_var_from_secret_scoped_to_module"
}

func (d *moduleEnvVarFromSecretScopedToModuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Module Inputs (Env Vars) --- Use this data source to access information about an existing Module Env Var (From Secret Scoped To Module) in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedId + "Module Env Var (From Secret Scoped To Module).",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedName1 + "Module Env Var (From Secret Scoped To Module). " + DescSharedName2,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedModuleId1 + "Module Env Var (From Secret Scoped To Module)" + DescSharedModuleId2,
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedSecretType,
			},
			"secret_scoped_to_module_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the Secret Scoped To Module to take as input.",
			},
		},
	}
}

func (d *moduleEnvVarFromSecretScopedToModuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleEnvVarFromSecretScopedToModuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleEnvVarFromSecretScopedToModuleEndpoint, data.ModuleId.ValueString(), data.SecretScopedToModuleId.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromSecretScopedToModuleDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromSecretScopedToModuleDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
