// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var moduleParamFromOutputDefaultError = fmt.Sprintf("snapcd_moduleParamFromOutput error")

var moduleParamFromOutputEndpoint = "/api/Definition/ModuleParamFromOutput"

var _ resource.Resource = (*moduleParamFromOutputResource)(nil)

func ModuleParamFromOutputResource() resource.Resource {
	return &moduleParamFromOutputResource{}
}

type moduleParamFromOutputResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *moduleParamFromOutputResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *moduleParamFromOutputResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_param_from_output"
}

type moduleParamFromOutputModel struct {
	Name          types.String `tfsdk:"name"`
	Id            types.String `tfsdk:"id"`
	OutputName    types.String `tfsdk:"output_name"`
	ModuleName    types.String `tfsdk:"module_name"`
	NamespaceName types.String `tfsdk:"namespace_name"`
	ModuleId      types.String `tfsdk:"module_id"`
}

func (r *moduleParamFromOutputResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"output_name": schema.StringAttribute{
				Required: true,
			},
			"module_name": schema.StringAttribute{
				Required: true,
			},
			"namespace_name": schema.StringAttribute{
				Required: true,
			},
			"module_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *moduleParamFromOutputResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data moduleParamFromOutputModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, err := r.client.Post(moduleParamFromOutputEndpoint, jsonMap)
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleParamFromOutputResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data moduleParamFromOutputModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, err := r.client.Get(fmt.Sprintf("%s/%s", moduleParamFromOutputEndpoint, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleParamFromOutputResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data moduleParamFromOutputModel
	var state moduleParamFromOutputModel

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
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, err := r.client.Put(fmt.Sprintf("%s/%s", moduleParamFromOutputEndpoint, state.Id.ValueString()), jsonMap)
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleParamFromOutputResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data moduleParamFromOutputModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.client.Delete(fmt.Sprintf("%s/%s", moduleParamFromOutputEndpoint, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *moduleParamFromOutputResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data moduleParamFromOutputModel

	result, err := r.client.Get(fmt.Sprintf("%s/%s", moduleParamFromOutputEndpoint, req.ID))
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromOutputDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
