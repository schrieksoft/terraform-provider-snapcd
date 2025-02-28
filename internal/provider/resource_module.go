package provider

import (
	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var moduleDefaultError = fmt.Sprintf("snapcd_module error")

var moduleEndpoint = "/api/Definition/Module"

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

type moduleModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	NamespaceId              types.String `tfsdk:"namespace_id"`
	RunnerPoolId             types.String `tfsdk:"runner_pool_id"`
	TargetRepoRevision       types.String `tfsdk:"target_repo_revision"`
	TargetRepoUrl            types.String `tfsdk:"target_repo_url"`
	TargetModuleRelativePath types.String `tfsdk:"target_module_relative_path"`
	ProviderCacheEnabled     types.Bool   `tfsdk:"provider_cache_enabled"`
	ModuleCacheEnabled       types.Bool   `tfsdk:"module_cache_enabled"`
	DependsOnModules         types.List   `tfsdk:"depends_on_modules"`
	SelectOn                 types.String `tfsdk:"select_on"`
	SelectStrategy           types.String `tfsdk:"select_strategy"`
	SelectedConsumerId       types.String `tfsdk:"selected_consumer_id"`
	InitBeforeHook           types.String `tfsdk:"init_before_hook"`
	InitAfterHook            types.String `tfsdk:"init_after_hook"`
	InitBackendArgs          types.String `tfsdk:"init_backend_args"`
	PlanBeforeHook           types.String `tfsdk:"plan_before_hook"`
	PlanAfterHook            types.String `tfsdk:"plan_after_hook"`
	ApplyBeforeHook          types.String `tfsdk:"apply_before_hook"`
	ApplyAfterHook           types.String `tfsdk:"apply_after_hook"`
	PlanDestroyBeforeHook    types.String `tfsdk:"plan_destroy_before_hook"`
	PlanDestroyAfterHook     types.String `tfsdk:"plan_destroy_after_hook"`
	DestroyBeforeHook        types.String `tfsdk:"destroy_before_hook"`
	DestroyAfterHook         types.String `tfsdk:"destroy_after_hook"`
	OutputBeforeHook         types.String `tfsdk:"output_before_hook"`
	OutputAfterHook          types.String `tfsdk:"output_after_hook"`
	Engine                   types.String `tfsdk:"engine"`
}

func (r *moduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"namespace_id": schema.StringAttribute{
				Required: true,
			},
			"runner_pool_id": schema.StringAttribute{
				Required: true,
			},
			"target_repo_revision": schema.StringAttribute{
				Optional: true,
			},
			"target_repo_url": schema.StringAttribute{
				Optional: true,
			},
			"target_module_relative_path": schema.StringAttribute{
				Required: true,
			},
			"provider_cache_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"module_cache_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"depends_on_modules": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"select_on": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("PoolId", "ConsumerId"),
				},
			},
			"select_strategy": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("FirstOf", "AnyOf"),
				},
			},
			"selected_consumer_id": schema.StringAttribute{
				Optional: true,
			},
			"init_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"init_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"init_backend_args": schema.StringAttribute{
				Optional: true,
			},
			"plan_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"plan_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"plan_destroy_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"plan_destroy_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"apply_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"apply_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"destroy_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"destroy_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"output_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"output_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"engine": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("OpenTofu", "Terraform"),
				},
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
	if httpError != nil && httpError.StatusCode == 442 {
		resp.Diagnostics.AddError(globalRoleAssignmentDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
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
	if httpError != nil && httpError.StatusCode == 441 {
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
	if httpError != nil && httpError.StatusCode == 441 {
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
