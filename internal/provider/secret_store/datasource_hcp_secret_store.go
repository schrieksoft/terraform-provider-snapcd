package secret_store

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*hcpSecretStoreDataSource)(nil)

func HcpSecretStoreDataSource() datasource.DataSource {
	return &hcpSecretStoreDataSource{}
}

type hcpSecretStoreDataSource struct {
	client *snapcd.Client
}

func (r *hcpSecretStoreDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *hcpSecretStoreDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hcp_secret_store"
}

func (d *hcpSecretStoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Secret Stores --- Use this data source to access information about an existing HCP Secret Store in Snap CD.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: DescName,
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: DescId,
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: DescProjectId,
				Computed:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: DescOrganizationId,
				Computed:            true,
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: DescAppName,
				Computed:            true,
			},
			"is_assigned_to_all_scopes": schema.BoolAttribute{
				MarkdownDescription: DescIsAssignedToAllScopes,
				Computed:            true,
			},
		},
	}
}

func (d *hcpSecretStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hcpSecretStoreModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/ByName/%s", hcpSecretStoreEndpoint, data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(hcpSecretStoreDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(hcpSecretStoreDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
