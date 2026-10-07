package identity

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var servicePrincipalDefaultError = fmt.Sprintf("snapcd_service_principal error")

var servicePrincipalEndpoint = "/ServicePrincipal"

var _ datasource.DataSource = (*servicePrincipalDataSource)(nil)

func ServicePrincipalDataSource() datasource.DataSource {
	return &servicePrincipalDataSource{}
}

// Client ID, display name and ID only: a data source exists to resolve an ID, and narrowing it
// to metadata means a principal who may only discover the Service Principal can still use it.
// The client ID is the bare name, without the organization prefix it is stored with.
type servicePrincipalModel struct {
	Id          types.String `tfsdk:"id"`
	ClientId    types.String `tfsdk:"client_id"`
	DisplayName types.String `tfsdk:"display_name"`
}

type servicePrincipalDataSource struct {
	client *snapcd.Client
}

func (r *servicePrincipalDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *servicePrincipalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_principal"
}

func (d *servicePrincipalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Identity Access Management --- Look up an existing Service Principal by client ID to obtain its ID. Returns client ID, display name and ID only." + "\n\n## Required permissions\n\n" + openapidocs.DataSourceMetadataPermissions["ServicePrincipal"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ServicePrincipalMetadataReadDto_Id,
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ServicePrincipalMetadataReadDto_ClientId,
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ServicePrincipalMetadataReadDto_DisplayName,
			},
		},
	}
}

func (d *servicePrincipalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data servicePrincipalModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/Metadata/ByClientId/%s", servicePrincipalEndpoint, data.ClientId.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
