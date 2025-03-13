package core

import (
	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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

type namespaceModel struct {
	Name                         types.String `tfsdk:"name"`
	Id                           types.String `tfsdk:"id"`
	StackId                      types.String `tfsdk:"stack_id"`
	DefaultInitBeforeHook        types.String `tfsdk:"default_init_before_hook"`
	DefaultInitAfterHook         types.String `tfsdk:"default_init_after_hook"`
	DefaultInitBackendArgs       types.String `tfsdk:"default_init_backend_args"`
	DefaultPlanBeforeHook        types.String `tfsdk:"default_plan_before_hook"`
	DefaultPlanAfterHook         types.String `tfsdk:"default_plan_after_hook"`
	DefaultApplyBeforeHook       types.String `tfsdk:"default_apply_before_hook"`
	DefaultApplyAfterHook        types.String `tfsdk:"default_apply_after_hook"`
	DefaultPlanDestroyBeforeHook types.String `tfsdk:"default_plan_destroy_before_hook"`
	DefaultPlanDestroyAfterHook  types.String `tfsdk:"default_plan_destroy_after_hook"`
	DefaultDestroyBeforeHook     types.String `tfsdk:"default_destroy_before_hook"`
	DefaultDestroyAfterHook      types.String `tfsdk:"default_destroy_after_hook"`
	DefaultOutputBeforeHook      types.String `tfsdk:"default_output_before_hook"`
	DefaultOutputAfterHook       types.String `tfsdk:"default_output_after_hook"`
	DefaultEngine                types.String `tfsdk:"default_engine"`
	DefaultTargetRepoRevision    types.String `tfsdk:"default_target_repo_revision"`
	DefaultTargetRepoUrl         types.String `tfsdk:"default_target_repo_url"`
	DefaultProviderCacheEnabled  types.Bool   `tfsdk:"default_provider_cache_enabled"`
	DefaultModuleCacheEnabled    types.Bool   `tfsdk:"default_module_cache_enabled"`
	DefaultSelectOn              types.String `tfsdk:"default_select_on"`
	DefaultSelectStrategy        types.String `tfsdk:"default_select_strategy"`
	DefaultOutputSecretStoreId   types.String `tfsdk:"default_output_secret_store_id"`
}

func (r *namespaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (r *namespaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"stack_id": schema.StringAttribute{
				Required: true,
			},
			// "auto_refresh": schema.BoolAttribute{
			// 	Optional: true,
			// 	Computed: true,
			// 	Default:  booldefault.StaticBool(false),
			// },
			// "current_state_refresh_interval": schema.StringAttribute{
			// 	Optional: true,
			// 	Computed: true,
			// 	Default:  stringdefault.StaticString("120m"),
			// },
			"default_init_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_init_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_init_backend_args": schema.StringAttribute{
				Optional: true,
			},
			"default_plan_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_plan_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_plan_destroy_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_plan_destroy_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_apply_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_apply_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_destroy_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_destroy_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_output_before_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_output_after_hook": schema.StringAttribute{
				Optional: true,
			},
			"default_engine": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("OpenTofu", "Terraform"),
				},
			},

			"default_target_repo_revision": schema.StringAttribute{
				Optional: true,
			},

			"default_target_repo_url": schema.StringAttribute{
				Optional: true,
			},

			"default_provider_cache_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"default_module_cache_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"default_select_on": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("PoolId", "ConsumerId"),
				},
			},
			"default_select_strategy": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("FirstOf", "AnyOf"),
				},
			},
			"default_output_secret_store_id": schema.StringAttribute{
				Optional: true,
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
