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

var resourceRoleAssignmentDefaultError = fmt.Sprintf("snapcd_resource_role_assignment error")

var resourceRoleAssignmentEndpoint = "/api/Identity/ResourceRoleAssignment"

var _ resource.Resource = (*resourceRoleAssignmentResource)(nil)

func ResourceRoleAssignmentResource() resource.Resource {
	return &resourceRoleAssignmentResource{}
}

type resourceRoleAssignmentResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *resourceRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_role_assignment"
}

type resourceRoleAssignmentModel struct {
	Id                     types.String `tfsdk:"id"`
	ResourceId             types.String `tfsdk:"resource_id"`
	PrincipalId            types.String `tfsdk:"principal_id"`
	PrincipalDiscriminator types.String `tfsdk:"principal_discriminator"`
	ResourceDiscriminator  types.String `tfsdk:"resource_discriminator"`
	RoleName               types.String `tfsdk:"role_name"`
}

func (r *resourceRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.StringAttribute{
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
			"resource_discriminator": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Stack", "Namespace", "Module", "RunnerPool"),
				},
			},
			"role_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Owner", "Contributor", "Reader"),
				},
			},
		},
	}
}

func (r *resourceRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resourceRoleAssignmentModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(resourceRoleAssignmentEndpoint, jsonMap)
	err = httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})

	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resourceRoleAssignmentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", resourceRoleAssignmentEndpoint, data.Id.ValueString()))
	err := httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resourceRoleAssignmentModel
	var state resourceRoleAssignmentModel

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
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", resourceRoleAssignmentEndpoint, state.Id.ValueString()), jsonMap)
	err = httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resourceRoleAssignmentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", resourceRoleAssignmentEndpoint, data.Id.ValueString()))
	err := httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *resourceRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data resourceRoleAssignmentModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", resourceRoleAssignmentEndpoint, req.ID))
	err := httpError.Error
	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(resourceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
