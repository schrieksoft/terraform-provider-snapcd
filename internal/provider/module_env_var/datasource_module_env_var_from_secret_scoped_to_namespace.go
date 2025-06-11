package module_env_var

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*moduleEnvVarFromSecretScopedToNamespaceDataSource)(nil)

func ModuleEnvVarFromSecretScopedToNamespaceDataSource() datasource.DataSource {
	return &moduleEnvVarFromSecretScopedToNamespaceDataSource{}
}

type moduleEnvVarFromSecretScopedToNamespaceDataSource struct {
	client *snapcd.Client
}

func (r *moduleEnvVarFromSecretScopedToNamespaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleEnvVarFromSecretScopedToNamespaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_env_var_from_secret_scoped_to_namespace"
}

func (d *moduleEnvVarFromSecretScopedToNamespaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Module Inputs (Env Vars) --- Use this data source to access information about an existing Module Env Var (From Secret Scoped To Namespace) in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedId + "Module Env Var (From Secret Scoped To Namespace).",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedName1 + "Module Env Var (From Secret Scoped To Namespace). " + DescSharedName2,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedModuleId1 + "Module Env Var (From Secret Scoped To Namespace)" + DescSharedModuleId2,
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedSecretType,
			},
			"secret_scoped_to_namespace_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the Secret Scoped To Namespace to take as input.",
			},
		},
	}
}

func (d *moduleEnvVarFromSecretScopedToNamespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleEnvVarFromSecretScopedToNamespaceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleEnvVarFromSecretScopedToNamespaceEndpoint, data.ModuleId.ValueString(), data.SecretScopedToNamespaceId.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromSecretScopedToNamespaceDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromSecretScopedToNamespaceDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
