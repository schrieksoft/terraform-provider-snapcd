// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*namespaceDataSource)(nil)

func NamespaceDataSource() datasource.DataSource {
	return &namespaceDataSource{}
}

type namespaceDataSource struct {
	client *snapcd.Client
}

func (r *namespaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (d *namespaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"stack_id": schema.StringAttribute{
				Required: true,
			},
			"default_init_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_init_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_init_backend_args": schema.StringAttribute{
				Computed: true,
			},
			"default_plan_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_plan_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_plan_destroy_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_plan_destroy_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_apply_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_apply_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_destroy_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_destroy_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_output_before_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_output_after_hook": schema.StringAttribute{
				Computed: true,
			},
			"default_engine": schema.StringAttribute{
				Computed: true,
			},

			"default_target_repo_revision": schema.StringAttribute{
				Computed: true,
			},

			"default_target_repo_url": schema.StringAttribute{
				Computed: true,
			},

			"default_provider_cache_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"default_module_cache_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"default_select_on": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("PoolId", "ConsumerId"),
				},
			},
			"default_select_strategy": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("FirstOf", "AnyOf"),
				},
			},
		},
	}
}

func (d *namespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Get(fmt.Sprintf("%s/%s/%s", namespaceEndpoint, data.StackId.ValueString(), data.Name.ValueString()))

	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
