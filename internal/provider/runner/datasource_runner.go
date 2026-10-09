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

var _ datasource.DataSource = (*runnerDataSource)(nil)

func RunnerDataSource() datasource.DataSource {
	return &runnerDataSource{}
}

type runnerDataSource struct {
	client *snapcd.Client
}

func (r *runnerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *runnerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runner"
}

func (d *runnerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Runners --- Look up an existing Runner by name to obtain its ID. Returns name and ID only." + "\n\n## Required permissions\n\n" + openapidocs.DataSourceMetadataPermissions["Runner"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.RunnerMetadataReadDto_Id,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.RunnerMetadataReadDto_Name,
			},
			"is_disabled": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.RunnerMetadataReadDto_IsDisabled,
			},
		},
	}
}

func (d *runnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data runnerMetadataModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/Metadata/ByName/%s", runnerEndpoint, data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(runnerDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(runnerDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Name and ID only: a data source exists to resolve an ID, and narrowing it to metadata
// means a principal who may only discover the resource can still use it.
type runnerMetadataModel struct {
	Name       types.String `tfsdk:"name"`
	Id         types.String `tfsdk:"id"`
	IsDisabled types.Bool   `tfsdk:"is_disabled"`
}
