package secret

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*azureKeyVaultSecretScopedToModuleDataSource)(nil)

func AzureKeyVaultSecretScopedToModuleDataSource() datasource.DataSource {
	return &azureKeyVaultSecretScopedToModuleDataSource{}
}

type azureKeyVaultSecretScopedToModuleDataSource struct {
	client *snapcd.Client
}

func (r *azureKeyVaultSecretScopedToModuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *azureKeyVaultSecretScopedToModuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_key_vault_secret_scoped_to_module"
}

func (d *azureKeyVaultSecretScopedToModuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Secrets --- Use this data source to access information about an existing Azure Key Vault Secret (Scoped to Module) in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescName,
			},
			"remote_secret_name": schema.StringAttribute{
				Computed:    true,
				Description: DescRemoteName,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescModuleId,
			},
			"secret_store_id": schema.StringAttribute{
				Computed:    true,
				Description: DescSecretStoreId,
			},
		},
	}
}

func (d *azureKeyVaultSecretScopedToModuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data azureKeyVaultSecretScopedToModuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/ByName/%s", azureKeyVaultSecretScopedToModuleEndpoint, data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToModuleDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToModuleDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
