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
	_ datasource.DataSource              = &namespaceEnvVarFromSecretScopedToStackDataSource{}
	_ datasource.DataSourceWithConfigure = &namespaceEnvVarFromSecretScopedToStackDataSource{}
)

func NamespaceEnvVarFromSecretScopedToStackDataSource() datasource.DataSource {
	return &namespaceEnvVarFromSecretScopedToStackDataSource{}
}

type namespaceEnvVarFromSecretScopedToStackDataSource struct {
	client *snapcd.Client
}

func (d *namespaceEnvVarFromSecretScopedToStackDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_env_var_from_secret_scoped_to_stack"
}

func (d *namespaceEnvVarFromSecretScopedToStackDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Namespace Inputs (Env Vars) --- Retrieves a Namespace Env Var (From Secret Scoped To Stack) from Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedId + "Namespace Env Var (From Secret Scoped To Stack).",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedName1 + "Namespace Env Var (From Secret Scoped To Stack). " + DescSharedName2,
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedSecretType,
			},
			"secret_scoped_to_stack_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the secret scoped to stack to use for this namespace env var.",
			},
			"usage_mode": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedUsage,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedNamespaceId1 + "Namespace Env Var (From Secret Scoped To Stack)" + DescSharedNamespaceId2,
			},
		},
	}
}

func (d *namespaceEnvVarFromSecretScopedToStackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespaceEnvVarFromSecretScopedToStackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceEnvVarFromSecretScopedToStackModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", namespaceEnvVarFromSecretScopedToStackEndpoint, data.NamespaceId.ValueString(), data.SecretScopedToStackId.ValueString()))

	if httpError != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromSecretScopedToStackDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromSecretScopedToStackDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
