package secret

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var secretScopedToModuleDefaultError = fmt.Sprintf("snapcd_secret_scoped_to_module error")

var secretScopedToModuleEndpoint = "/SecretScopedToModule"

type secretScopedToModuleModel struct {
	Name     types.String `tfsdk:"name"`
	Id       types.String `tfsdk:"id"`
	ModuleId types.String `tfsdk:"module_id"`
}

var _ datasource.DataSource = (*secretScopedToModuleDataSource)(nil)

func SecretScopedToModuleDataSource() datasource.DataSource {
	return &secretScopedToModuleDataSource{}
}

type secretScopedToModuleDataSource struct {
	client *snapcd.Client
}

func (r *secretScopedToModuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *secretScopedToModuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret_scoped_to_module"
}

func (d *secretScopedToModuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Secrets --- Use this data source to access information about an existing Secret (Scoped to Module) in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescName,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescModuleId,
			},
		},
	}
}

func (d *secretScopedToModuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data secretScopedToModuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/ByName/%s", secretScopedToModuleEndpoint, data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(secretScopedToModuleDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(secretScopedToModuleDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
