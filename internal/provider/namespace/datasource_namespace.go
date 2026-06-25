package namespace

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

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
		MarkdownDescription: "Namespaces --- Use this data source to access information about an existing Namespace in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescNamespaceId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceName,
			},
			"stack_id": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceStackId,
			},
			"default_clean_init_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: DescNamespaceCleanInitEnabled,
			},
			"default_drift_check_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: DescNamespaceDefaultDriftCheckEnabled,
			},
			"default_drift_check_interval_minutes": schema.Int64Attribute{
				Computed:    true,
				Description: DescNamespaceDefaultDriftCheckIntervalMinutes,
			},
			"default_engine": schema.StringAttribute{
				Computed:    true,
				Description: DescNamespaceDefaultEngine,
			},
			"trigger_behaviour_on_modified": schema.StringAttribute{
				Computed:    true,
				Description: DescNamespaceTriggerBehaviourOnModified,
			},
			"default_apply_approval_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: DescNamespaceDefaultApplyApprovalThreshold,
			},
			"default_destroy_approval_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: DescNamespaceDefaultDestroyApprovalThreshold,
			},
			"default_approval_timeout_minutes": schema.Int64Attribute{
				Computed:    true,
				Description: DescNamespaceDefaultApprovalTimeoutMinutes,
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

	result, httpError := d.client.Get(fmt.Sprintf("%s/%s/%s", namespaceEndpoint, data.StackId.ValueString(), data.Name.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

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
