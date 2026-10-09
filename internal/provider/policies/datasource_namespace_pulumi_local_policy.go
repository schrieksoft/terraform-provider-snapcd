package policies

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

var _ datasource.DataSource = (*namespacePulumiLocalPolicyDataSource)(nil)

func NamespacePulumiLocalPolicyDataSource() datasource.DataSource {
	return &namespacePulumiLocalPolicyDataSource{}
}

type namespacePulumiLocalPolicyDataSource struct {
	client *snapcd.Client
}

func (r *namespacePulumiLocalPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespacePulumiLocalPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_pulumi_local_policy"
}

func (d *namespacePulumiLocalPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Policies --- Use this data source to access information about an existing Namespace Pulumi Local Policy in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourceMetadataPermissions["NamespacePulumiLocalPolicy"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespacePulumiLocalPolicyMetadataReadDto_Id,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespacePulumiLocalPolicyMetadataReadDto_NamespaceId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespacePulumiLocalPolicyMetadataReadDto_Name,
			},
		},
	}
}

func (d *namespacePulumiLocalPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespacePulumiLocalPolicyMetadataModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/Metadata/%s/%s", namespacePulumiLocalPolicyEndpoint, data.NamespaceId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(namespacePulumiLocalPolicyDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespacePulumiLocalPolicyDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Metadata only: a data source exists to resolve an ID, and narrowing it means a
// principal who may only discover the resource can still use it.
type namespacePulumiLocalPolicyMetadataModel struct {
	Id          types.String `tfsdk:"id"`
	NamespaceId types.String `tfsdk:"namespace_id"`
	Name        types.String `tfsdk:"name"`
}
