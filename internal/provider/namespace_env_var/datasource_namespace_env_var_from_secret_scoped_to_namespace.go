package namespace_env_var

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &namespaceEnvVarFromSecretScopedToNamespaceDataSource{}
	_ datasource.DataSourceWithConfigure = &namespaceEnvVarFromSecretScopedToNamespaceDataSource{}
)

func NamespaceEnvVarFromSecretScopedToNamespaceDataSource() datasource.DataSource {
	return &namespaceEnvVarFromSecretScopedToNamespaceDataSource{}
}

type namespaceEnvVarFromSecretScopedToNamespaceDataSource struct {
	client *snapcd.Client
}

func (d *namespaceEnvVarFromSecretScopedToNamespaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_env_var_from_secret_scoped_to_namespace"
}

func (d *namespaceEnvVarFromSecretScopedToNamespaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Namespace Inputs (Env Vars) --- Retrieves a Namespace Env Var (From Secret Scoped To Namespace) from Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedId + "Namespace Env Var (From Secret Scoped To Namespace).",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescSharedName1 + "Namespace Env Var (From Secret Scoped To Namespace). " + DescSharedName2,
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedSecretType,
			},
			"secret_scoped_to_namespace_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the secret scoped to namespace to use for this namespace env var.",
			},
			"usage_mode": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedUsage,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedNamespaceId1 + "Namespace Env Var (From Secret Scoped To Namespace)" + DescSharedNamespaceId2,
			},
		},
	}
}

func (d *namespaceEnvVarFromSecretScopedToNamespaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*snapcd.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *snapcd.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *namespaceEnvVarFromSecretScopedToNamespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceEnvVarFromSecretScopedToNamespaceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", namespaceEnvVarFromSecretScopedToNamespaceEndpoint, data.NamespaceId.ValueString(), data.Name.ValueString()))
	
	if httpError != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromSecretScopedToNamespaceDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromSecretScopedToNamespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
