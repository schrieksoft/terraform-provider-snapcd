package core

import (
	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var namespaceDefaultError = fmt.Sprintf("snapcd_namespace error")

var namespaceEndpoint = "/api/Namespace"

var _ resource.Resource = (*namespaceResource)(nil)

func NamespaceResource() resource.Resource {
	return &namespaceResource{}
}

type namespaceResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *namespaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ! Category: Namespace
type namespaceModel struct {
	Name                            types.String `tfsdk:"name"`
	Id                              types.String `tfsdk:"id"`
	StackId                         types.String `tfsdk:"stack_id"`
	DefaultInitBeforeHook           types.String `tfsdk:"default_init_before_hook"`
	DefaultInitAfterHook            types.String `tfsdk:"default_init_after_hook"`
	DefaultInitBackendArgs          types.String `tfsdk:"default_init_backend_args"`
	DefaultPlanBeforeHook           types.String `tfsdk:"default_plan_before_hook"`
	DefaultPlanAfterHook            types.String `tfsdk:"default_plan_after_hook"`
	DefaultApplyBeforeHook          types.String `tfsdk:"default_apply_before_hook"`
	DefaultApplyAfterHook           types.String `tfsdk:"default_apply_after_hook"`
	DefaultPlanDestroyBeforeHook    types.String `tfsdk:"default_plan_destroy_before_hook"`
	DefaultPlanDestroyAfterHook     types.String `tfsdk:"default_plan_destroy_after_hook"`
	DefaultDestroyBeforeHook        types.String `tfsdk:"default_destroy_before_hook"`
	DefaultDestroyAfterHook         types.String `tfsdk:"default_destroy_after_hook"`
	DefaultOutputBeforeHook         types.String `tfsdk:"default_output_before_hook"`
	DefaultOutputAfterHook          types.String `tfsdk:"default_output_after_hook"`
	DefaultEngine                   types.String `tfsdk:"default_engine"`
	DefaultOutputSecretStoreId      types.String `tfsdk:"default_output_secret_store_id"`
	DefaultApplyApprovalThreshold   types.Number `tfsdk:"default_apply_approval_threshold"`
	DefaultDestroyApprovalThreshold types.Number `tfsdk:"default_destroy_approval_threshold"`

	TriggerBehaviourOnModified types.String `tfsdk:"trigger_behaviour_on_modified"`
}

const (
	DescNamespaceDefault = "All modules in this Namespace will use this value, unless explicitly overriden on the Module itself."

	DescNamespaceId      = "Unique ID of the Namespace"
	DescNamespaceName    = "Name of the Namespace. Must be unique in combination with `stack_id`."
	DescNamespaceStackId = "ID of the Namespace's parent Stack."

	DescNamespaceDefaultInitBackendArgs          = DescSharedInitBackedArgs + DescNamespaceDefault
	DescNamespaceDefaultInitBeforeHook           = DescSharedInitBeforeHook + DescNamespaceDefault
	DescNamespaceDefaultInitAfterHook            = DescSharedInitAfterHook + DescNamespaceDefault
	DescNamespaceDefaultPlanBeforeHook           = DescSharedPlanBeforeHook + DescNamespaceDefault
	DescNamespaceDefaultPlanAfterHook            = DescSharedPlanAfterHook + DescNamespaceDefault
	DescNamespaceDefaultPlanDestroyBeforeHook    = DescSharedPlanDestroyBeforeHook + DescNamespaceDefault
	DescNamespaceDefaultPlanDestroyAfterHook     = DescSharedPlanDestroyAfterHook + DescNamespaceDefault
	DescNamespaceDefaultApplyBeforeHook          = DescSharedApplyBeforeHook + DescNamespaceDefault
	DescNamespaceDefaultApplyAfterHook           = DescSharedApplyAfterHook + DescNamespaceDefault
	DescNamespaceDefaultDestroyBeforeHook        = DescSharedDestroyBeforeHook + DescNamespaceDefault
	DescNamespaceDefaultDestroyAfterHook         = DescSharedDestroyAfterHook + DescNamespaceDefault
	DescNamespaceDefaultOutputBeforeHook         = DescSharedOutputBeforeHook + DescNamespaceDefault
	DescNamespaceDefaultOutputAfterHook          = DescSharedOutputAfterHook + DescNamespaceDefault
	DescNamespaceDefaultEngine                   = DescSharedEngine + DescNamespaceDefault
	DescNamespaceDefaultOutputSecretStoreId      = DescSharedOutputSecretStoreId + DescNamespaceDefault
	DescNamespaceDefaultApplyApprovalThreshold   = DescSharedApplyApprovalThreshold + DescNamespaceDefault + DescZeroThreshold
	DescNamespaceDefaultDestroyApprovalThreshold = DescSharedDestroyApprovalThreshold + DescNamespaceDefault + DescZeroThreshold

	DescNamespaceTriggerBehaviourOnModified = "Behaviour with respect to applying modules within the Namespace if any of the fields on the Namespace resource (or any of its Param, Env Var or Extra File resources) has changed. Must be one of 'TriggerAllImmediately' or 'DoNotTrigger'. Setting to 'TriggerAllImmediately' will trigger *all* Modules within the Stack to run an apply Job simultaneously. Setting to 'DoNotTrigger' will do nothing. The default (and recommended) setting is 'DoNotTrigger'."
)

func (r *namespaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (r *namespaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Namespaces --- Manages a Namespace in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"default_init_backend_args": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultInitBackendArgs,
			},
			"default_init_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultInitBeforeHook,
			},
			"default_init_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultInitAfterHook,
			},
			"default_plan_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultPlanBeforeHook,
			},
			"default_plan_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultPlanAfterHook,
			},
			"default_plan_destroy_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultPlanDestroyBeforeHook,
			},
			"default_plan_destroy_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultPlanDestroyAfterHook,
			},
			"default_apply_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultApplyBeforeHook,
			},
			"default_apply_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultApplyAfterHook,
			},
			"default_destroy_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultDestroyBeforeHook,
			},
			"default_destroy_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultDestroyAfterHook,
			},
			"default_output_before_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultOutputBeforeHook,
			},
			"default_output_after_hook": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultOutputAfterHook,
			},
			"default_engine": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("OpenTofu", "Terraform"),
				},
				Description: DescNamespaceDefaultEngine,
			},
			"default_output_secret_store_id": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultOutputSecretStoreId,
			},
			"default_apply_approval_threshold": schema.NumberAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultApplyApprovalThreshold,
			},

			"default_destroy_approval_threshold": schema.NumberAttribute{
				Optional:    true,
				Description: DescNamespaceDefaultDestroyApprovalThreshold,
			},
			"trigger_behaviour_on_modified": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: DescStackTriggerBehaviourOnModified,
				Validators: []validator.String{
					stringvalidator.OneOf("DoNotTrigger", "TriggerAllImmediately"),
				},
				Default: stringdefault.StaticString("DoNotTrigger"),
			},
		},
	}
}

func (r *namespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data namespaceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(namespaceEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(namespaceDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data namespaceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data namespaceModel
	var state namespaceModel

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
		resp.Diagnostics.AddError(namespaceDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", namespaceEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data namespaceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", namespaceEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *namespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data namespaceModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
