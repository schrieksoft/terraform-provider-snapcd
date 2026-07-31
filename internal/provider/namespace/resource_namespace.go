package namespace

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

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

var namespaceEndpoint = "/Namespace"

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
	Name                             types.String `tfsdk:"name"`
	Id                               types.String `tfsdk:"id"`
	StackId                          types.String `tfsdk:"stack_id"`
	DefaultEngine                    types.String `tfsdk:"default_engine"`
	DefaultApplyApprovalThreshold    types.Int64  `tfsdk:"default_apply_approval_threshold"`
	DefaultDestroyApprovalThreshold  types.Int64  `tfsdk:"default_destroy_approval_threshold"`
	DefaultApprovalTimeoutMinutes    types.Int64  `tfsdk:"default_approval_timeout_minutes"`
	DefaultCleanInitEnabled          types.Bool   `tfsdk:"default_clean_init_enabled"`
	DefaultTriggerPathFilterEnabled  types.Bool   `tfsdk:"default_trigger_path_filter_enabled"`
	DefaultDriftCheckEnabled         types.Bool   `tfsdk:"default_drift_check_enabled"`
	DefaultDriftCheckIntervalMinutes types.Int64  `tfsdk:"default_drift_check_interval_minutes"`
	TriggerBehaviourOnModified       types.String `tfsdk:"trigger_behaviour_on_modified"`
}

const (
	DescNamespaceDefault = "All modules in this Namespace will use this value, unless explicitly overriden on the Module itself."
)

func (r *namespaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (r *namespaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Namespaces --- Manages a Namespace in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.ResourcePermissions["Namespace"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: openapidocs.NamespaceReadDto_Id,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceCreateDto_Name,
			},
			"stack_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceCreateDto_StackId,
			},
			"default_clean_init_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: openapidocs.NamespaceCreateDto_DefaultCleanInitEnabled,
			},
			"default_trigger_path_filter_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: openapidocs.NamespaceCreateDto_DefaultTriggerPathFilterEnabled,
			},
			"default_drift_check_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: openapidocs.NamespaceCreateDto_DefaultDriftCheckEnabled,
			},
			"default_drift_check_interval_minutes": schema.Int64Attribute{
				Optional:    true,
				Description: openapidocs.NamespaceCreateDto_DefaultDriftCheckIntervalMinutes,
			},

			"default_engine": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.StateManagementEngineValues...),
				},
				Description: openapidocs.NamespaceCreateDto_DefaultEngine,
			},
			"default_apply_approval_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: openapidocs.NamespaceCreateDto_DefaultApplyApprovalThreshold,
			},

			"default_destroy_approval_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: openapidocs.NamespaceCreateDto_DefaultDestroyApprovalThreshold,
			},

			"default_approval_timeout_minutes": schema.Int64Attribute{
				Optional:    true,
				Description: openapidocs.NamespaceCreateDto_DefaultApprovalTimeoutMinutes,
			},

			"trigger_behaviour_on_modified": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.NamespaceCreateDto_TriggerBehaviourOnModified,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.NamespaceTriggerBehaviourValues...),
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
