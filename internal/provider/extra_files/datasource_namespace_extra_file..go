package extra_files

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*namespaceExtraFileDataSource)(nil)

func NamespaceExtraFileDataSource() datasource.DataSource {
	return &namespaceExtraFileDataSource{}
}

type namespaceExtraFileDataSource struct {
	client *snapcd.Client
}

func (r *namespaceExtraFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespaceExtraFileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_extra_file"
}

func (d *namespaceExtraFileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Extra Files --- Use this data source to access information about an existing Namespace Extra File in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescNamespaceExtraFileId,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceExtraFileNamespaceId,
			},
			"file_name": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceExtraFileFilename,
			},
			"contents": schema.StringAttribute{
				Computed:    true,
				Description: DescNamespaceExtraFileContents,
			},
			"overwrite": schema.BoolAttribute{
				Computed:    true,
				Description: DescNamespaceExtraFileOverwrite,
			},
		},
	}
}

func (d *namespaceExtraFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceExtraFileModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", namespaceExtraFileEndpoint, data.NamespaceId.ValueString(), data.FileName.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(namespaceExtraFileDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceExtraFileDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
