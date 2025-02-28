package provider

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

var groupMemberDefaultError = fmt.Sprintf("snapcd_group_member error")

var groupMemberEndpoint = "/api/Identity/GroupMember"

var _ resource.Resource = (*groupMemberResource)(nil)

func GroupMemberResource() resource.Resource {
	return &groupMemberResource{}
}

type groupMemberResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *groupMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_member"
}

type groupMemberModel struct {
	Id                     types.String `tfsdk:"id"`
	GroupId                types.String `tfsdk:"group_id"`
	PrincipalId            types.String `tfsdk:"principal_id"`
	PrincipalDiscriminator types.String `tfsdk:"principal_discriminator"`
}

func (r *groupMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required: true,
			},
			"principal_id": schema.StringAttribute{
				Required: true,
			},
			"principal_discriminator": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("User", "ServicePrincipal", "Group"),
				},
			},
		},
	}
}

func (r *groupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data groupMemberModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(groupMemberEndpoint, jsonMap)
	err = httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})

	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data groupMemberModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", groupMemberEndpoint, data.Id.ValueString()))
	err := httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data groupMemberModel
	var state groupMemberModel

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
		resp.Diagnostics.AddError(groupMemberDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", groupMemberEndpoint, state.Id.ValueString()), jsonMap)
	err = httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data groupMemberModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", groupMemberEndpoint, data.Id.ValueString()))
	err := httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *groupMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data groupMemberModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", groupMemberEndpoint, req.ID))
	err := httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(groupMemberDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
