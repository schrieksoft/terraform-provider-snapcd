package namespace_input

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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &namespaceInputFromSecretDataSource{}
	_ datasource.DataSourceWithConfigure = &namespaceInputFromSecretDataSource{}
)

func NamespaceInputFromSecretDataSource() datasource.DataSource {
	return &namespaceInputFromSecretDataSource{}
}

type namespaceInputFromSecretDataSource struct {
	client *snapcd.Client
}

func (d *namespaceInputFromSecretDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_input_from_secret"
}

func (d *namespaceInputFromSecretDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Namespace Inputs --- Retrieves a Namespace Input (From Secret) from Snap CD.` + "\n\n## Required permissions\n\n" + openapidocs.DataSourceMetadataPermissions["NamespaceInputFromSecret"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.NamespaceInputMetadataReadDto_Id,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceInputMetadataReadDto_Name,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceInputMetadataReadDto_NamespaceId,
			},
			"input_kind": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceInputFromSecretReadDto_InputKind,
			},
		},
	}
}

func (d *namespaceInputFromSecretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*snapcd.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *snapcd.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *namespaceInputFromSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceInputFromSecretMetadataModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := d.client.Get(fmt.Sprintf("%s/Metadata/%s/%s", namespaceInputFromSecretEndpoint, data.NamespaceId.ValueString(), data.Name.ValueString()))

	if httpError != nil {
		resp.Diagnostics.AddError(namespaceInputFromSecretDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceInputFromSecretDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Metadata only: a data source exists to resolve an ID, and narrowing it means a
// principal who may only discover the resource can still use it.
type namespaceInputFromSecretMetadataModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	NamespaceId types.String `tfsdk:"namespace_id"`
	InputKind   types.String `tfsdk:"input_kind"`
}
