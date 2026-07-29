package module

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

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
	Engine                             types.String `tfsdk:"engine"`
	TriggerOnDefinitionChanged         types.Bool   `tfsdk:"trigger_on_definition_changed"`
	TriggerOnUpstreamOutputChanged     types.Bool   `tfsdk:"trigger_on_upstream_output_changed"`
	TriggerOnSourceChanged             types.Bool   `tfsdk:"trigger_on_source_changed"`
	TriggerOnSourceChangedNotification types.Bool   `tfsdk:"trigger_on_source_changed_notification"`
	ApplyApprovalThreshold             types.Int64  `tfsdk:"apply_approval_threshold"`
	DestroyApprovalThreshold           types.Int64  `tfsdk:"destroy_approval_threshold"`
	ApprovalTimeoutMinutes             types.Int64  `tfsdk:"approval_timeout_minutes"`
	CleanInitEnabled                   types.Bool   `tfsdk:"clean_init_enabled"`
	DriftCheckEnabled                  types.Bool   `tfsdk:"drift_check_enabled"`
	DriftCheckIntervalMinutes          types.Int64  `tfsdk:"drift_check_interval_minutes"`
	IgnoreNamespaceExtraFiles          types.Bool   `tfsdk:"ignore_namespace_extra_files"`
	IgnoreNamespaceFlags               types.Bool   `tfsdk:"ignore_namespace_flags"`
	IgnoreNamespaceHooks               types.Bool   `tfsdk:"ignore_namespace_hooks"`
	WaitForApplyDependencies           types.String `tfsdk:"wait_for_apply_dependencies"`
	WaitForDestroyDependencies         types.String `tfsdk:"wait_for_destroy_dependencies"`
}

func (r *moduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Modules --- Manages a Module in Snap CD.` + "\n\n## Required permissions\n\n" + openapidocs.ResourcePermissions["Module"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: openapidocs.ModuleReadDto_Id,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleCreateDto_Name,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleCreateDto_NamespaceId,
			},
			"runner_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleCreateDto_RunnerId,
			},
			"source_revision": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleCreateDto_SourceRevision,
			},
			"source_url": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ModuleCreateDto_SourceUrl,
			},
			"source_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					// Intentionally narrower than the spec enum: the other source types are internal.
					stringvalidator.OneOf("Git", "Registry"),
				},
				Default:     stringdefault.StaticString("Git"),
				Description: openapidocs.ModuleCreateDto_SourceType,
			},
			"source_revision_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.SourceRevisionTypeValues...),
				},
				Default:     stringdefault.StaticString("Default"),
				Description: openapidocs.ModuleCreateDto_SourceRevisionType,
			},
			"source_subdirectory": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: openapidocs.ModuleCreateDto_SourceSubdirectory,
			},
			"runner_instance_name": schema.StringAttribute{
				Optional:    true,
				Description: openapidocs.ModuleCreateDto_RunnerInstanceName,
			},
			"clean_init_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: openapidocs.ModuleCreateDto_CleanInitEnabled,
			},
			"drift_check_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: openapidocs.ModuleCreateDto_DriftCheckEnabled,
			},
			"drift_check_interval_minutes": schema.Int64Attribute{
				Optional:    true,
				Description: openapidocs.ModuleCreateDto_DriftCheckIntervalMinutes,
			},
			"ignore_namespace_extra_files": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.ModuleCreateDto_IgnoreNamespaceExtraFiles,
				Default:     booldefault.StaticBool(false),
			},
			"ignore_namespace_flags": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.ModuleCreateDto_IgnoreNamespaceFlags,
				Default:     booldefault.StaticBool(false),
			},
			"ignore_namespace_hooks": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.ModuleCreateDto_IgnoreNamespaceHooks,
				Default:     booldefault.StaticBool(false),
			},
			"engine": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.StateManagementEngineValues...),
				},
				Description: openapidocs.ModuleCreateDto_Engine,
			},
			"apply_approval_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: openapidocs.ModuleCreateDto_ApplyApprovalThreshold,
			},

			"destroy_approval_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: openapidocs.ModuleCreateDto_DestroyApprovalThreshold,
			},
			"approval_timeout_minutes": schema.Int64Attribute{
				Optional:    true,
				Description: openapidocs.ModuleCreateDto_ApprovalTimeoutMinutes,
			},

			"trigger_on_definition_changed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.ModuleCreateDto_TriggerOnDefinitionChanged,
				Default:     booldefault.StaticBool(true),
			},
			"trigger_on_upstream_output_changed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.ModuleCreateDto_TriggerOnUpstreamOutputChanged,
				Default:     booldefault.StaticBool(true),
			},
			"trigger_on_source_changed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.ModuleCreateDto_TriggerOnSourceChanged,
				Default:     booldefault.StaticBool(true),
			},
			"trigger_on_source_changed_notification": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.ModuleCreateDto_TriggerOnSourceChangedNotification,
				Default:     booldefault.StaticBool(false),
			},
			"wait_for_apply_dependencies": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.WaitForApplyDependenciesValues...),
				},
				Default:     stringdefault.StaticString("OnFirstApply"),
				Description: openapidocs.ModuleCreateDto_WaitForApplyDependencies,
			},
			"wait_for_destroy_dependencies": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.WaitForDestroyDependenciesValues...),
				},
				Default:     stringdefault.StaticString("Always"),
				Description: openapidocs.ModuleCreateDto_WaitForDestroyDependencies,
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
