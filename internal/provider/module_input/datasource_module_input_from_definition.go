package module_input

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*moduleInputFromDefinitionDataSource)(nil)

func ModuleInputFromDefinitionDataSource() datasource.DataSource {
	return &moduleInputFromDefinitionDataSource{}
}

type moduleInputFromDefinitionDataSource struct {
	client *snapcd.Client
}

func (r *moduleInputFromDefinitionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleInputFromDefinitionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_input_from_definition"
}

func (d *moduleInputFromDefinitionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Module Inputs --- Use this data source to access information about an existing Module Input (From Definition) in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["ModuleInputFromDefinition"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleInputFromDefinitionReadDto_Id,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleInputFromDefinitionReadDto_Name,
			},
			"definition_name": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleInputFromDefinitionReadDto_DefinitionName,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleInputFromDefinitionReadDto_ModuleId,
			},
			"input_kind": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleInputFromDefinitionReadDto_InputKind,
			},
		},
	}
}

func (d *moduleInputFromDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleInputFromDefinitionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleInputFromDefinitionEndpoint, data.ModuleId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleInputFromDefinitionDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleInputFromDefinitionDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
