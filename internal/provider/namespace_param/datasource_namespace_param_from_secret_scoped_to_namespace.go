package namespace_param

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
	_ datasource.DataSource              = &namespaceParamFromSecretScopedToNamespaceDataSource{}
	_ datasource.DataSourceWithConfigure = &namespaceParamFromSecretScopedToNamespaceDataSource{}
)

func NamespaceParamFromSecretScopedToNamespaceDataSource() datasource.DataSource {
	return &namespaceParamFromSecretScopedToNamespaceDataSource{}
}

type namespaceParamFromSecretScopedToNamespaceDataSource struct {
	client *snapcd.Client
}

func (d *namespaceParamFromSecretScopedToNamespaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_param_from_secret_scoped_to_namespace"
}

func (d *namespaceParamFromSecretScopedToNamespaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Namespace Inputs (Parameters) --- Retrieves a Namespace Param (From Secret Scoped To Namespace) from Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedId + "Namespace Param (From Secret Scoped To Namespace).",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedName1 + "Namespace Param (From Secret Scoped To Namespace). " + DescSharedName2,
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedSecretType,
			},
			"secret_scoped_to_namespace_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the secret scoped to namespace to use for this namespace param.",
			},
			"usage_mode": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedUsage,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedNamespaceId1 + "Namespace Param (From Secret Scoped To Namespace)" + DescSharedNamespaceId2,
			},
		},
	}
}

func (d *namespaceParamFromSecretScopedToNamespaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespaceParamFromSecretScopedToNamespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceParamFromSecretScopedToNamespaceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", namespaceParamFromSecretScopedToNamespaceEndpoint, data.NamespaceId.ValueString(), data.SecretScopedToNamespaceId.ValueString()))

	if httpError != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToNamespaceDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToNamespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
