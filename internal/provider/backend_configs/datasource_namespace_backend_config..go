package backend_configs

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*namespaceBackendConfigDataSource)(nil)

func NamespaceBackendConfigDataSource() datasource.DataSource {
	return &namespaceBackendConfigDataSource{}
}

type namespaceBackendConfigDataSource struct {
	client *snapcd.Client
}

func (r *namespaceBackendConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespaceBackendConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_backend_config"
}

func (d *namespaceBackendConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Backend Configs --- Use this data source to access information about an existing Namespace Backend Config in Snap CD.",
		DeprecationMessage:  "Use snapcd_namespace_terraform_array_flag with Flag='BackendConfig' and Task='Init' instead.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescNamespaceBackendConfigId,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceBackendConfigNamespaceId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceBackendConfigName,
			},
			"value": schema.StringAttribute{
				Computed:    true,
				Description: DescNamespaceBackendConfigValue,
			},
		},
	}
}

func (d *namespaceBackendConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceBackendConfigModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", namespaceBackendConfigEndpoint, data.NamespaceId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(namespaceBackendConfigDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceBackendConfigDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
