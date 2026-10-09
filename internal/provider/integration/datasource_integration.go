package integration

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"
)

var integrationDataSourceError = "snapcd_integration data source error"

var _ datasource.DataSource = (*integrationDataSource)(nil)

func IntegrationDataSource() datasource.DataSource {
	return &integrationDataSource{}
}

type integrationDataSource struct {
	client *snapcd.Client
}

// The integration carries a write-only secret, so it is created/managed in the SnapCd UI, not by Terraform.
// This data source looks one up by name so its id can be referenced by snapcd_integration_*_supply /
// snapcd_integration_role_assignment / snapcd_integration_event. The connection (credentials) is never
// exposed here — it is not returned by the API for this read and never lands in Terraform state.
type integrationDataSourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	IntegrationType        types.String `tfsdk:"integration_type"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	IsSuppliedToAllModules types.Bool   `tfsdk:"is_supplied_to_all_modules"`
}

func (d *integrationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*snapcd.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *snapcd.Client, got: %T.", req.ProviderData))
		return
	}
	d.client = client
}

func (d *integrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (d *integrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Integrations --- Look up an existing Integration by name to obtain its ID. Returns name, ID and type only." + "\n\n## Required permissions\n\n" + openapidocs.DataSourceMetadataPermissions["Integration"],
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Computed: true, Description: openapidocs.IntegrationMetadataReadDto_Id},
			"name":             schema.StringAttribute{Required: true, Description: openapidocs.IntegrationMetadataReadDto_Name},
			"integration_type": schema.StringAttribute{Computed: true, Description: openapidocs.IntegrationMetadataReadDto_IntegrationType},
		},
	}
}

func (d *integrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data integrationMetadataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/Metadata/ByName/%s", integrationEndpoint, data.Name.ValueString()))
	if httpError != nil {
		resp.Diagnostics.AddError(integrationDataSourceError, "Error calling GET: "+httpError.Error.Error())
		return
	}

	if err := utils.JsonToPlan(result, &data); err != nil {
		resp.Diagnostics.AddError(integrationDataSourceError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Name and ID only: a data source exists to resolve an ID, and narrowing it to metadata
// means a principal who may only discover the resource can still use it.
type integrationMetadataModel struct {
	Name            types.String `tfsdk:"name"`
	Id              types.String `tfsdk:"id"`
	IntegrationType types.String `tfsdk:"integration_type"`
}
