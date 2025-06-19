package secret_store

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &hcpSecretStoreDataSource{}

func HCPSecretStoreDataSource() datasource.DataSource {
	return &hcpSecretStoreDataSource{}
}

var hcpSecretStoreDataSourceDefaultError = fmt.Sprintf("snapcd_hcp_secret_store data source error")

// hcpSecretStoreDataSource defines the data source implementation.
type hcpSecretStoreDataSource struct {
	client *snapcd.Client
}

// hcpSecretStoreDataSourceModel describes the data source data model.
type hcpSecretStoreDataSourceModel struct {
	Name                  types.String `tfsdk:"name"`
	Id                    types.String `tfsdk:"id"`
	ProjectId             types.String `tfsdk:"project_id"`
	OrganizationId        types.String `tfsdk:"organization_id"`
	AppName               types.String `tfsdk:"app_name"`
	IsAssignedToAllScopes types.Bool   `tfsdk:"is_assigned_to_all_scopes"`
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

func (d *hcpSecretStoreDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *hcpSecretStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data hcpSecretStoreDataSourceModel

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
		resp.Diagnostics.AddError(hcpSecretStoreDataSourceDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(hcpSecretStoreDataSourceDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
