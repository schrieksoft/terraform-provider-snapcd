package module

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
		MarkdownDescription: "Modules --- Use this data source to access information about an existing Module in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.DataSourcePermissions["Module"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_Id,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleReadDto_Name,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleReadDto_NamespaceId,
			},
			"runner_id": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_RunnerId,
			},
			"source_revision": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_SourceRevision,
			},
			"source_url": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_SourceUrl,
			},
			"source_subdirectory": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_SourceSubdirectory,
			},
			"source_type": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					// Intentionally narrower than the spec enum: the other source types are internal.
					stringvalidator.OneOf("Git", "Registry"),
				},
				Description: openapidocs.ModuleReadDto_SourceType,
			},
			"source_revision_type": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					// Intentionally narrower than the spec enum: the data source resolves Default only.
					stringvalidator.OneOf("Default"),
				},
				Description: openapidocs.ModuleReadDto_SourceRevisionType,
			},
			"runner_instance_name": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_RunnerInstanceName,
			},
			"clean_init_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_CleanInitEnabled,
			},
			"ignore_namespace_extra_files": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_IgnoreNamespaceExtraFiles,
			},
			"ignore_namespace_flags": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_IgnoreNamespaceFlags,
			},
			"ignore_namespace_hooks": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_IgnoreNamespaceHooks,
			},
			"drift_check_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_DriftCheckEnabled,
			},
			"drift_check_interval_minutes": schema.Int64Attribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_DriftCheckIntervalMinutes,
			},

			"engine": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.StateManagementEngineValues...),
				},
				Description: openapidocs.ModuleReadDto_Engine,
			},

			"trigger_on_definition_changed": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_TriggerOnDefinitionChanged,
			},
			"trigger_on_upstream_output_changed": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_TriggerOnUpstreamOutputChanged,
			},
			"trigger_on_source_changed": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_TriggerOnSourceChanged,
			},
			"trigger_on_source_changed_notification": schema.BoolAttribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_TriggerOnSourceChangedNotification,
			},
			"apply_approval_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_ApplyApprovalThreshold,
			},
			"destroy_approval_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_DestroyApprovalThreshold,
			},
			"approval_timeout_minutes": schema.Int64Attribute{
				Computed:    true,
				Description: openapidocs.ModuleReadDto_ApprovalTimeoutMinutes,
			},
			"wait_for_apply_dependencies": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.WaitForApplyDependenciesValues...),
				},
				Description: openapidocs.ModuleReadDto_WaitForApplyDependencies,
			},
			"wait_for_destroy_dependencies": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.WaitForDestroyDependenciesValues...),
				},
				Description: openapidocs.ModuleReadDto_WaitForDestroyDependencies,
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
