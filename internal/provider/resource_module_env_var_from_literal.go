// Copyright (c) HashiCorp, Inc.

package provider

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

var moduleEnvVarFromLiteralDefaultError = fmt.Sprintf("snapcd_moduleEnvVarFromLiteral error")

var moduleEnvVarFromLiteralEndpoint = "/api/Definition/ModuleEnvVarFromLiteral"

var _ resource.Resource = (*moduleEnvVarFromLiteralResource)(nil)

func ModuleEnvVarFromLiteralResource() resource.Resource {
	return &moduleEnvVarFromLiteralResource{}
}

type moduleEnvVarFromLiteralResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *moduleEnvVarFromLiteralResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *moduleEnvVarFromLiteralResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_env_var_from_literal"
}

type moduleEnvVarFromLiteralModel struct {
	Name         types.String `tfsdk:"name"`
	Id           types.String `tfsdk:"id"`
	LiteralValue types.String `tfsdk:"literal_value"`
	Type         types.String `tfsdk:"type"`
	ModuleId     types.String `tfsdk:"module_id"`
}

func (r *moduleEnvVarFromLiteralResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"literal_value": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("String", "NotString", "Number", "Bool", "Tuple", "Object")},
				Default: stringdefault.StaticString("String"),
			},
			"module_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *moduleEnvVarFromLiteralResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data moduleEnvVarFromLiteralModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, err := r.client.Post(moduleEnvVarFromLiteralEndpoint, jsonMap)
	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleEnvVarFromLiteralResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data moduleEnvVarFromLiteralModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, err := r.client.Get(fmt.Sprintf("%s/%s", moduleEnvVarFromLiteralEndpoint, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleEnvVarFromLiteralResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data moduleEnvVarFromLiteralModel
	var state moduleEnvVarFromLiteralModel

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
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, err := r.client.Put(fmt.Sprintf("%s/%s", moduleEnvVarFromLiteralEndpoint, state.Id.ValueString()), jsonMap)
	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleEnvVarFromLiteralResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data moduleEnvVarFromLiteralModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.client.Delete(fmt.Sprintf("%s/%s", moduleEnvVarFromLiteralEndpoint, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *moduleEnvVarFromLiteralResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data moduleEnvVarFromLiteralModel

	result, err := r.client.Get(fmt.Sprintf("%s/%s", moduleEnvVarFromLiteralEndpoint, req.ID))
	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
