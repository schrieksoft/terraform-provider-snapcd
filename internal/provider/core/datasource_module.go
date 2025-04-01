package core

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*moduleDataSource)(nil)

func ModuleDataSource() datasource.DataSource {
	return &moduleDataSource{}
}

type moduleDataSource struct {
	client *snapcd.Client
}

func (r *moduleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *moduleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module"
}

func (d *moduleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"namespace_id": schema.StringAttribute{
				Required: true,
			},
			"runner_pool_id": schema.StringAttribute{
				Computed: true,
			},
			"source_revision": schema.StringAttribute{
				Computed: true,
			},
			"source_url": schema.StringAttribute{
				Computed: true,
			},
			"source_subdirectory": schema.StringAttribute{
				Computed: true,
			},
			"depends_on_modules": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"source_type": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Git"),
				},
			},
			"source_revision_type": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Default"),
				},
			},
			"select_on": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("None", "PoolId", "ConsumerId"),
				},
			},
			"select_strategy": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("FirstOf", "AnyOf"),
				},
			},
			"selected_consumer_id": schema.StringAttribute{
				Computed: true,
			},
			"init_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"init_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"init_backend_args": schema.StringAttribute{
				Computed: true,
			},
			"plan_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"plan_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"plan_destroy_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"plan_destroy_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"apply_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"apply_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"destroy_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"destroy_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"output_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"output_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"engine": schema.StringAttribute{
				Computed: true,
			},
			"output_secret_store_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *moduleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data moduleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", moduleEndpoint, data.NamespaceId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
