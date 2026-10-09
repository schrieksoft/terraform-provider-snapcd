package runner

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*sourceRefresherPreselectionDataSource)(nil)

func SourceRefresherPreselectionDataSource() datasource.DataSource {
	return &sourceRefresherPreselectionDataSource{}
}

type sourceRefresherPreselectionDataSource struct {
	client *snapcd.Client
}

func (r *sourceRefresherPreselectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *sourceRefresherPreselectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_refresher_preselection"
}

func (d *sourceRefresherPreselectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Runners --- Use this data source to access information about an existing Source Refresher Preselection in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourceMetadataPermissions["SourceRefresherPreselection"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.SourceRefresherPreselectionMetadataReadDto_Id,
			},
			"source_url": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.SourceRefresherPreselectionReadDto_SourceUrl,
			},
			"runner_id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.SourceRefresherPreselectionMetadataReadDto_RunnerId,
			},
		},
	}
}

func (d *sourceRefresherPreselectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data sourceRefresherPreselectionMetadataModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/Metadata/BySourceUrl/%s", sourceRefresherPreselectionEndpoint, data.SourceUrl.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Metadata only: a data source exists to resolve an ID, and narrowing it means a
// principal who may only discover the resource can still use it.
type sourceRefresherPreselectionMetadataModel struct {
	Id        types.String `tfsdk:"id"`
	SourceUrl types.String `tfsdk:"source_url"`
	RunnerId  types.String `tfsdk:"runner_id"`
}
