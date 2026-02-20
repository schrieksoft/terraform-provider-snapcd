package core

import (
	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var moduleDefaultError = fmt.Sprintf("snapcd_module error")

var moduleEndpoint = "/Module"

var _ resource.Resource = (*moduleResource)(nil)

func ModuleResource() resource.Resource {
	return &moduleResource{}
}

type moduleResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *moduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *moduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module"
}

// ! Category: Module
type moduleModel struct {
	Id                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	NamespaceId                        types.String `tfsdk:"namespace_id"`
	RunnerId                           types.String `tfsdk:"runner_id"`
	SourceRevision                     types.String `tfsdk:"source_revision"`
	SourceUrl                          types.String `tfsdk:"source_url"`
	SourceSubdirectory                 types.String `tfsdk:"source_subdirectory"`
	SourceType                         types.String `tfsdk:"source_type"`
	SourceRevisionType                 types.String `tfsdk:"source_revision_type"`
	RunnerInstanceName                 types.String `tfsdk:"runner_instance_name"`
	InitBeforeHook                     types.String `tfsdk:"init_before_hook"`
	InitAfterHook                      types.String `tfsdk:"init_after_hook"`
	PlanBeforeHook                     types.String `tfsdk:"plan_before_hook"`
	PlanAfterHook                      types.String `tfsdk:"plan_after_hook"`
	ApplyBeforeHook                    types.String `tfsdk:"apply_before_hook"`
	ApplyAfterHook                     types.String `tfsdk:"apply_after_hook"`
	PlanDestroyBeforeHook              types.String `tfsdk:"plan_destroy_before_hook"`
	PlanDestroyAfterHook               types.String `tfsdk:"plan_destroy_after_hook"`
	DestroyBeforeHook                  types.String `tfsdk:"destroy_before_hook"`
	DestroyAfterHook                   types.String `tfsdk:"destroy_after_hook"`
	OutputBeforeHook                   types.String `tfsdk:"output_before_hook"`
	OutputAfterHook                    types.String `tfsdk:"output_after_hook"`
	ValidateBeforeHook                 types.String `tfsdk:"validate_before_hook"`
	ValidateAfterHook                  types.String `tfsdk:"validate_after_hook"`
	Engine                             types.String `tfsdk:"engine"`
	OutputSecretStoreId                types.String `tfsdk:"output_secret_store_id"`
	TriggerOnDefinitionChanged         types.Bool   `tfsdk:"trigger_on_definition_changed"`
	TriggerOnUpstreamOutputChanged     types.Bool   `tfsdk:"trigger_on_upstream_output_changed"`
	TriggerOnSourceChanged             types.Bool   `tfsdk:"trigger_on_source_changed"`
	TriggerOnSourceChangedNotification types.Bool   `tfsdk:"trigger_on_source_changed_notification"`
	ApplyApprovalThreshold             types.Int64  `tfsdk:"apply_approval_threshold"`
	DestroyApprovalThreshold           types.Int64  `tfsdk:"destroy_approval_threshold"`
	ApprovalTimeoutMinutes             types.Int64  `tfsdk:"approval_timeout_minutes"`
	AutoUpgradeEnabled                 types.Bool   `tfsdk:"auto_upgrade_enabled"`
	AutoReconfigureEnabled             types.Bool   `tfsdk:"auto_reconfigure_enabled"`
	AutoMigrateEnabled                 types.Bool   `tfsdk:"auto_migrate_enabled"`
	CleanInitEnabled                   types.Bool   `tfsdk:"clean_init_enabled"`
	IgnoreNamespaceBackendConfigs      types.Bool   `tfsdk:"ignore_namespace_backend_configs"`
	IgnoreNamespaceExtraFiles          types.Bool   `tfsdk:"ignore_namespace_extra_files"`
	WaitForApplyDependencies           types.String `tfsdk:"wait_for_apply_dependencies"`
	WaitForDestroyDependencies         types.String `tfsdk:"wait_for_destroy_dependencies"`
}

const (
	DescModuleOverride = "Setting this will override any default value set on the Module's parent Namespace."

	DescModuleId                            = "Unique ID of the Module."
	DescModuleName                          = "Name of the Module. Must be unique in combination with `namespace_id`."
	DescModuleNamespaceId                   = "ID of the Module's parent Namespace."
	DescModuleRunnerId                      = "ID of the Runner that will receive the instructions when triggering a deployment on this Module."
	DescModuleSourceRevision                = "Remote revision (e.g. version number, branch, commit or tag) where the source module code is found."
	DescModuleSourceUrl                     = "Remote URL where the source module code is found."
	DescModuleSourceSubdirectory            = "Subdirectory where the source module code is found."
	DescModuleDependsOnModules              = "A list on Snap CD Modules that this Module depends on. Note that Snap CD will automatically discover depedencies based on the Module using as inputs the outputs from another Module, so use `depends_on_modules` where you want to explicitly establish a dependency where outputs are not referenced as inputs."
	DescModuleSourceType                    = "The type of remote module store that the source module code should be retrieved from. Must be one of 'Git' or 'Registry'"
	DescModuleSourceRevisionType            = "How Snap CD should interpret the `source_revision` field. Must be one of 'Default' or 'SemanticVersionRange'. Setting to 'Default' means Snap CD will interpret the revision type based on the source type (for example, for a 'Git' source type it will automatically figure out whether the `source_revision` refers to a branch, tag or commit). Setting to 'SemanticVersionRange' means that Snap CD will resolve the revision to a semantic version line `vX.Y.Z` (alternatively witout the 'v' prefix of that is how your semantic version are tagged, i.e. 'X.Y.Z'). It will take the highest version within the major or minor version range that you specify. For example, specify `v2.20.*` or `v2.*`. You can also specify a specific semantic version here, e.g. `v2.20.7`. In that case the behaviour is the same as with when using 'Default', except that only valid semantic versions are accepted. NOTE that 'SemanticVersionRange' is currently only supported in combination with the 'Git' `source_type`."
	DescModuleRunnerInstanceName            = "Name a specific runner instance to select (should unique identify the the instance). Use this if you have enabled multiple instances on your runner, but want all jobs for this Module to go to a specific instance."
	DescModuleInitBackendArgs               = DescSharedInitBackedArgs + DescModuleOverride
	DescModuleInitBeforeHook                = DescSharedInitBeforeHook + DescModuleOverride
	DescModuleInitAfterHook                 = DescSharedInitAfterHook + DescModuleOverride
	DescModulePlanBeforeHook                = DescSharedPlanBeforeHook + DescModuleOverride
	DescModulePlanAfterHook                 = DescSharedPlanAfterHook + DescModuleOverride
	DescModulePlanDestroyBeforeHook         = DescSharedPlanDestroyBeforeHook + DescModuleOverride
	DescModulePlanDestroyAfterHook          = DescSharedPlanDestroyAfterHook + DescModuleOverride
	DescModuleApplyBeforeHook               = DescSharedApplyBeforeHook + DescModuleOverride
	DescModuleApplyAfterHook                = DescSharedApplyAfterHook + DescModuleOverride
	DescModuleDestroyBeforeHook             = DescSharedDestroyBeforeHook + DescModuleOverride
	DescModuleDestroyAfterHook              = DescSharedDestroyAfterHook + DescModuleOverride
	DescModuleOutputBeforeHook              = DescSharedOutputBeforeHook + DescModuleOverride
	DescModuleOutputAfterHook               = DescSharedOutputAfterHook + DescModuleOverride
	DescModuleValidateBeforeHook            = DescSharedValidateBeforeHook + DescModuleOverride
	DescModuleValidateAfterHook             = DescSharedValidateAfterHook + DescModuleOverride
	DescModuleEngine                        = DescSharedEngine + DescModuleOverride
	DescModuleOutputSecretStoreId           = DescSharedOutputSecretStoreId + DescModuleOverride
	DescModuleApplyApprovalThreshold        = DescSharedApplyApprovalThreshold + DescModuleOverride + DescZeroThreshold
	DescModuleDestroyApprovalThreshold      = DescSharedDestroyApprovalThreshold + DescModuleOverride + DescZeroThreshold
	DescModuleApprovalTimeoutMinutes        = DescSharedApprovalTimeoutMinutes + DescModuleOverride + DescZeroTimeout
	DescModuleAutoUpgradeEnabled            = DescSharedAutoUpgradeEnabled + DescModuleOverride
	DescModuleAutoReconfigureEnabled        = DescSharedAutoReconfigureEnabled + DescModuleOverride
	DescModuleAutoMigrateEnabled            = DescSharedAutoMigrateEnabled + DescModuleOverride
	DescModuleCleanInitEnabled              = DescSharedCleanInitEnabled + DescModuleOverride
	DescModuleIgnoreNamespaceBackendConfigs = "If this is set to true, any Backend Configs that have been set on Namespace level will not be used on this specific Module."
	DescModuleIgnoreNamespaceExtraFiles     = "If this is set to true, any Extra Files that have been set on Namespace level will not be used on this specific Module."
	DescTriggerOnSourceChanged             = "Defaults to 'true'. If 'true', the Module will automatically be applied when the source it is referencing has changed. For example, if tracking a Git branch: a new commit would constitute a change."
	DescTriggerOnSourceChangedNotification = "Defaults to 'false'. If 'true', the Module will automatically be applied when the 'api/Hooks/SourceChanged' endpoint is called for this Module. Use this if you want to use external tooling to inform Snap CD that a source has been changed. Consider setting `trigger_on_definition_changed` to 'false' when setting `trigger_on_definition_changed_hook` to 'true'"
	DescTriggerOnUpstreamOutputChanged     = "Defaults to 'true'. If 'true', the Module will automatically be applied when any Outputs from other Modules that it references as Inputs (Param or Env Var) have changed."
	DescTriggerOnDefinitionChanged         = "Defaults to 'true'. If 'true', the Module will automatically be applied when its definition changes. A definition change results from fields on the Module itself, on any of its Inputs (Param or Env Var) or Extra Files being altered. So too changes to its Namespace (including Inputs and Extra Files) or Stack. Note however that Namespace and Stack changes are not notified by default. This behaviour can be changed in `snapcd_namespace` and `snapcd_stack` resource definitions."
	DescWaitForApplyDependencies           = "Defaults to 'OnFirstApply'. Controls when the Module should wait for dependencies during apply operations. Valid values are 'Always', 'Never', or 'OnFirstApply'. 'Always' means the Module will always wait for Modules its depends on to reach the 'Applied' state before applying. 'Never' means dependencies are ignored. 'OnFirstApply' means the Module will wait for dependencies only on its first apply."
	DescWaitForDestroyDependencies         = "Defaults to 'Always'. Controls when the Module should wait for dependencies during destroy operations. Valid values are 'Always' or 'Never'. 'Always' means the Module will always wait Modules that depend on it to reach the 'Destroyed' state before destroying. 'Never' means dependencies are ignored."
)

func (r *moduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Modules --- Manages a Module in Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"runner_id": schema.StringAttribute{
				Required:    true,
				Description: DescModuleRunnerId,
			},
			"source_revision": schema.StringAttribute{
				Required:    true,
				Description: DescModuleSourceRevision,
			},
			"source_url": schema.StringAttribute{
				Required:    true,
				Description: DescModuleSourceUrl,
			},
			"source_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Git", "Registry"),
				},
				Default:     stringdefault.StaticString("Git"),
				Description: DescModuleSourceType,
			},
			"source_revision_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Default", "SemanticVersionRange"),
				},
				Default:     stringdefault.StaticString("Default"),
				Description: DescModuleSourceRevisionType,
			},
			"source_subdirectory": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: DescModuleSourceSubdirectory,
			},
			"runner_instance_name": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleRunnerInstanceName,
			},
			"init_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleInitBeforeHook,
			},
			"init_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleInitAfterHook,
			},
			"plan_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModulePlanBeforeHook,
			},
			"plan_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModulePlanAfterHook,
			},
			"plan_destroy_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleDestroyBeforeHook,
			},
			"plan_destroy_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleDestroyAfterHook,
			},
			"apply_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleApplyBeforeHook,
			},
			"apply_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleApplyAfterHook,
			},
			"destroy_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleDestroyBeforeHook,
			},
			"destroy_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleDestroyAfterHook,
			},
			"output_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleOutputBeforeHook,
			},
			"output_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleOutputAfterHook,
			},
			"validate_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleValidateBeforeHook,
			},
			"validate_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleValidateAfterHook,
			},

			"auto_upgrade_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: DescModuleAutoUpgradeEnabled,
			},
			"auto_reconfigure_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: DescModuleAutoReconfigureEnabled,
			},
			"auto_migrate_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: DescModuleAutoMigrateEnabled,
			},
			"clean_init_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: DescModuleCleanInitEnabled,
			},
			"ignore_namespace_extra_files": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: DescModuleIgnoreNamespaceExtraFiles,
				Default:     booldefault.StaticBool(false),
			},
			"ignore_namespace_backend_configs": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: DescModuleIgnoreNamespaceBackendConfigs,
				Default:     booldefault.StaticBool(false),
			},
			"engine": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("OpenTofu", "Terraform", "Pulumi"),
				},
				Description: DescModuleEngine,
			},
			"output_secret_store_id": schema.StringAttribute{
				Optional:    true,
				Description: DescModuleOutputSecretStoreId,
			},
			"apply_approval_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: DescModuleApplyApprovalThreshold,
			},

			"destroy_approval_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: DescModuleDestroyApprovalThreshold,
			},
			"approval_timeout_minutes": schema.Int64Attribute{
				Optional:    true,
				Description: DescModuleApprovalTimeoutMinutes,
			},

			"trigger_on_definition_changed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: DescTriggerOnDefinitionChanged,
				Default:     booldefault.StaticBool(true),
			},
			"trigger_on_upstream_output_changed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: DescTriggerOnUpstreamOutputChanged,
				Default:     booldefault.StaticBool(true),
			},
			"trigger_on_source_changed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: DescTriggerOnSourceChanged,
				Default:     booldefault.StaticBool(true),
			},
			"trigger_on_source_changed_notification": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: DescTriggerOnSourceChangedNotification,
				Default:     booldefault.StaticBool(false),
			},
			"wait_for_apply_dependencies": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Always", "Never", "OnFirstApply"),
				},
				Default:     stringdefault.StaticString("OnFirstApply"),
				Description: DescWaitForApplyDependencies,
			},
			"wait_for_destroy_dependencies": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Always", "Never"),
				},
				Default:     stringdefault.StaticString("Always"),
				Description: DescWaitForDestroyDependencies,
			},
		},
	}
}

func (r *moduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data moduleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(moduleEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(moduleDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data moduleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", moduleEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		// Resource was not found, so remove it from state
		resp.State.RemoveResource(ctx)
		return
	}
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data moduleModel
	var state moduleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", moduleEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data moduleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", moduleEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		// Resource was not found, so remove it from state
		resp.State.RemoveResource(ctx)
		return
	}
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *moduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data moduleModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", moduleEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
