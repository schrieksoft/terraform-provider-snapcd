package namespace_input

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*namespaceInputFromDefinitionDataSource)(nil)

func NamespaceInputFromDefinitionDataSource() datasource.DataSource {
	return &namespaceInputFromDefinitionDataSource{}
}

type namespaceInputFromDefinitionDataSource struct {
	client *snapcd.Client
}

func (r *namespaceInputFromDefinitionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespaceInputFromDefinitionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_input_from_definition"
}

func (d *namespaceInputFromDefinitionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Namespace Inputs --- Use this data source to access information about an existing Namesapce Param (From Definition) in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["NamespaceInputFromDefinition"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespaceInputFromDefinitionReadDto_Id,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceInputFromDefinitionReadDto_Name,
			},
			"definition_name": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespaceInputFromDefinitionReadDto_DefinitionName,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.DefinitionInputTypeValues...),
				},
			},
			"usage_mode": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespaceInputFromDefinitionReadDto_UsageMode,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceInputFromDefinitionReadDto_NamespaceId,
			},
			"input_kind": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceInputFromDefinitionReadDto_InputKind,
			},
		},
	}
}

func (d *namespaceInputFromDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceInputFromDefinitionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", namespaceInputFromDefinitionEndpoint, data.NamespaceId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(namespaceInputFromDefinitionDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceInputFromDefinitionDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
