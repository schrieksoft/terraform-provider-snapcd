package secret

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*simpleSecretScopedToStackDataSource)(nil)

func SimpleSecretScopedToStackDataSource() datasource.DataSource {
	return &simpleSecretScopedToStackDataSource{}
}

type simpleSecretScopedToStackDataSource struct {
	client *snapcd.Client
}

func (r *simpleSecretScopedToStackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *simpleSecretScopedToStackDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_simple_secret_scoped_to_stack"
}

func (d *simpleSecretScopedToStackDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Secrets --- Use this data source to access information about an existing Simple Secret (Scoped to Stack) in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: DescId,
			},
			"name": schema.StringAttribute{
				Required: true,
				Description: DescName,
			},
			"value": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				Description: DescValue,
			},
			"stack_id": schema.StringAttribute{
				Required: true,
				Description: DescStackId,
			},
			"secret_store_id": schema.StringAttribute{
				Computed: true,
				Description: DescSecretStoreId,
			},
		},
	}
}

func (d *simpleSecretScopedToStackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data simpleSecretScopedToStackModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/ByName/%s", simpleSecretScopedToStackEndpoint, data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(simpleSecretScopedToStackDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(simpleSecretScopedToStackDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
