package role_assignments

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

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

var agentRoleAssignmentDefaultError = fmt.Sprintf("snapcd_agent_role_assignment error")

var agentRoleAssignmentEndpoint = "/AgentRoleAssignment"

var _ resource.Resource = (*agentRoleAssignmentResource)(nil)

func AgentRoleAssignmentResource() resource.Resource {
	return &agentRoleAssignmentResource{}
}

type agentRoleAssignmentResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *agentRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *agentRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_role_assignment"
}

type agentRoleAssignmentModel struct {
	Id                     types.String `tfsdk:"id"`
	AgentId                types.String `tfsdk:"agent_id"`
	PrincipalId            types.String `tfsdk:"principal_id"`
	PrincipalDiscriminator types.String `tfsdk:"principal_discriminator"`
	RoleName               types.String `tfsdk:"role_name"`
}

func (r *agentRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Role Assignments --- Manages an Agent Role Assignment in Snap CD.` + "\n\n## Required permissions\n\n" + openapidocs.ResourcePermissions["AgentRoleAssignment"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: openapidocs.AgentRoleAssignmentReadDto_Id,
			},
			"agent_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.AgentRoleAssignmentReadDto_AgentId,
			},
			"principal_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.AgentRoleAssignmentReadDto_PrincipalId,
			},
			"principal_discriminator": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					// Intentionally narrower than the spec enum: `Base` is an internal discriminator.
					stringvalidator.OneOf("User", "ServicePrincipal", "Group"),
				},
				Description: openapidocs.AgentRoleAssignmentReadDto_PrincipalDiscriminator,
			},
			"role_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.AgentRoleValues...),
				},
				Description: openapidocs.AgentRoleAssignmentReadDto_RoleName,
			},
		},
	}
}

func (r *agentRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data agentRoleAssignmentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(agentRoleAssignmentEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *agentRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data agentRoleAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", agentRoleAssignmentEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
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
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *agentRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data agentRoleAssignmentModel
	var state agentRoleAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", agentRoleAssignmentEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *agentRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data agentRoleAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", agentRoleAssignmentEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
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
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *agentRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data agentRoleAssignmentModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", agentRoleAssignmentEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(agentRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
