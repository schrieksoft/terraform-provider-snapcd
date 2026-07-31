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

var _ datasource.DataSource = (*moduleAdditionalTriggerPathDataSource)(nil)

func ModuleAdditionalTriggerPathDataSource() datasource.DataSource {
	return &moduleAdditionalTriggerPathDataSource{}
}

type moduleAdditionalTriggerPathDataSource struct {
	client *snapcd.Client
}

func (r *moduleAdditionalTriggerPathDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleAdditionalTriggerPathDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_additional_trigger_path"
}

func (d *moduleAdditionalTriggerPathDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Trigger Paths --- Use this data source to access information about an existing Module Additional Trigger Path in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["ModuleAdditionalTriggerPath"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleAdditionalTriggerPathReadDto_Id,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleAdditionalTriggerPathReadDto_ModuleId,
			},
			"path": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleAdditionalTriggerPathReadDto_Path,
			},
		},
	}
}

func (d *moduleAdditionalTriggerPathDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleAdditionalTriggerPathModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleAdditionalTriggerPathEndpoint, data.ModuleId.ValueString(), data.Path.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleAdditionalTriggerPathDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleAdditionalTriggerPathDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
