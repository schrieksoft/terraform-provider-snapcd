package core

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
		MarkdownDescription: "Modules --- Use this data source to access information about an existing Module in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescModuleName,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescModuleNamespaceId,
			},
			"runner_pool_id": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleRunnerPoolId,
			},
			"source_revision": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleSourceRevision,
			},
			"source_url": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleSourceUrl,
			},
			"source_subdirectory": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleSourceSubdirectory,
			},
			"source_type": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Git", "Registry"),
				},
				Description: DescModuleSourceType,
			},
			"source_revision_type": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Default"),
				},
				Description: DescModuleSourceRevisionType,
			},
			"runner_self_declared_name": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleRunnerSelfDeclaredName,
			},
			"init_before_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleInitBeforeHook,
			},
			"init_after_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleInitAfterHook,
			},
			"plan_before_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModulePlanBeforeHook,
			},
			"plan_after_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModulePlanAfterHook,
			},
			"plan_destroy_before_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModulePlanDestroyBeforeHook,
			},
			"plan_destroy_after_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModulePlanDestroyAfterHook,
			},
			"apply_before_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleApplyBeforeHook,
			},
			"apply_after_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleApplyAfterHook,
			},
			"destroy_before_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleDestroyBeforeHook,
			},
			"destroy_after_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleDestroyAfterHook,
			},
			"auto_upgrade_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: DescModuleAutoUpgradeEnabled,
			},
			"auto_reconfigure_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: DescModuleAutoReconfigureEnabled,
			},
			"auto_migrate_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: DescModuleAutoMigrateEnabled,
			},
			"clean_init_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: DescModuleCleanInitEnabled,
			},
			"ignore_namespace_extra_files": schema.BoolAttribute{
				Computed:    true,
				Description: DescModuleIgnoreNamespaceExtraFiles,
			},
			"ignore_namespace_backend_configs": schema.BoolAttribute{
				Computed:    true,
				Description: DescModuleIgnoreNamespaceBackendConfigs,
			},

			"output_before_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleOutputBeforeHook,
			},
			"output_after_hook": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleOutputAfterHook,
			},
			"engine": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("OpenTofu", "Terraform"),
				},
				Description: DescModuleEngine,
			},
			"output_secret_store_id": schema.StringAttribute{
				Computed:    true,
				Description: DescModuleOutputSecretStoreId,
			},

			"trigger_on_definition_changed": schema.BoolAttribute{
				Computed:    true,
				Description: DescTriggerOnDefinitionChanged,
			},
			"trigger_on_upstream_output_changed": schema.BoolAttribute{
				Computed:    true,
				Description: DescTriggerOnUpstreamOutputChanged,
			},
			"trigger_on_source_changed": schema.BoolAttribute{
				Computed:    true,
				Description: DescTriggerOnSourceChanged,
			},
			"trigger_on_source_changed_notification": schema.BoolAttribute{
				Computed:    true,
				Description: DescTriggerOnSourceChangedNotification,
			},
			"apply_approval_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: DescModuleApplyApprovalThreshold,
			},
			"destroy_approval_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: DescModuleDestroyApprovalThreshold,
			},
			"approval_timeout_minutes": schema.Int64Attribute{
				Computed:    true,
				Description: DescModuleApprovalTimeoutMinutes,
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
