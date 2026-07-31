package trigger_paths

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*namespaceAdditionalTriggerPathDataSource)(nil)

func NamespaceAdditionalTriggerPathDataSource() datasource.DataSource {
	return &namespaceAdditionalTriggerPathDataSource{}
}

type namespaceAdditionalTriggerPathDataSource struct {
	client *snapcd.Client
}

func (r *namespaceAdditionalTriggerPathDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespaceAdditionalTriggerPathDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_additional_trigger_path"
}

func (d *namespaceAdditionalTriggerPathDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Trigger Paths --- Use this data source to access information about an existing Module Additional Trigger Path in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["NamespaceAdditionalTriggerPath"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespaceAdditionalTriggerPathReadDto_Id,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceAdditionalTriggerPathReadDto_NamespaceId,
			},
			"path": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceAdditionalTriggerPathReadDto_Path,
			},
		},
	}
}

func (d *namespaceAdditionalTriggerPathDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceAdditionalTriggerPathModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", namespaceAdditionalTriggerPathEndpoint, data.NamespaceId.ValueString(), data.Path.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(namespaceAdditionalTriggerPathDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceAdditionalTriggerPathDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
