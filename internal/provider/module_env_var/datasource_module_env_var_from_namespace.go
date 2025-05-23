package module_env_var

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*moduleEnvVarFromNamespaceDataSource)(nil)

func ModuleEnvVarFromNamespaceDataSource() datasource.DataSource {
	return &moduleEnvVarFromNamespaceDataSource{}
}

type moduleEnvVarFromNamespaceDataSource struct {
	client *snapcd.Client
}

func (r *moduleEnvVarFromNamespaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleEnvVarFromNamespaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_env_var_from_namespace"
}

func (d *moduleEnvVarFromNamespaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Module Inputs (Env Vars) --- Use this data source to access information about an existing Module Env Var (From Namespace) in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedId + "Module Env Var (From Namespace).",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescSharedName1 + "Module Env Var (From Namespace). " + DescSharedName2,
			},
			"reference_name": schema.StringAttribute{
				Computed:    true,
				Description: DescSharedReferenceName,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedModuleId1 + "Module Env Var (From Namespace)" + DescSharedModuleId2,
			},
		},
	}
}

func (d *moduleEnvVarFromNamespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleEnvVarFromNamespaceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleEnvVarFromNamespaceEndpoint, data.ModuleId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromNamespaceDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromNamespaceDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
