package extra_files

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*moduleExtraFileDataSource)(nil)

func ModuleExtraFileDataSource() datasource.DataSource {
	return &moduleExtraFileDataSource{}
}

type moduleExtraFileDataSource struct {
	client *snapcd.Client
}

func (r *moduleExtraFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleExtraFileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_extra_file"
}

func (d *moduleExtraFileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Extra Files --- Use this data source to access information about an existing Module Extra File in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleExtraFileId,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescModuleExtraFileModuleId,
			},
			"file_name": schema.StringAttribute{
				Required:    true,
				Description: DescModuleExtraFileFilename,
			},
			"contents": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleExtraFileContents,
			},
			"overwrite": schema.BoolAttribute{
				Optional:    true,
				Description: DescModuleExtraFileOverwrite,
			},
		},
	}
}

func (d *moduleExtraFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleExtraFileModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleExtraFileEndpoint, data.ModuleId.ValueString(), data.FileName.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleExtraFileDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleExtraFileDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
