package integration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"
)

var roleAssignmentDefaultError = "snapcd_integration_role_assignment error"
var roleAssignmentEndpoint = "/IntegrationRoleAssignment"

var _ resource.Resource = (*integrationRoleAssignmentResource)(nil)

func IntegrationRoleAssignmentResource() resource.Resource {
	return &integrationRoleAssignmentResource{}
}

type integrationRoleAssignmentResource struct {
	client *snapcd.Client
}

func (r *integrationRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*snapcd.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *snapcd.Client, got: %T.", req.ProviderData))
		return
	}
	r.client = client
}

func (r *integrationRoleAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_role_assignment"
}

// ! Category: Integration
type integrationRoleAssignmentModel struct {
	Id                     types.String `tfsdk:"id"`
	IntegrationId          types.String `tfsdk:"integration_id"`
	PrincipalId            types.String `tfsdk:"principal_id"`
	PrincipalDiscriminator types.String `tfsdk:"principal_discriminator"`
	RoleName               types.String `tfsdk:"role_name"`
}

func (r *integrationRoleAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	requiresReplace := []planmodifier.String{stringplanmodifier.RequiresReplace()}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Integrations --- Grants an integration role to a principal on a specific integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "Unique ID of the role assignment.",
			},
			"integration_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: requiresReplace,
				Description:   "ID of the integration the role is granted on.",
			},
			"principal_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: requiresReplace,
				Description:   "ID of the principal (user / group / service principal).",
			},
			"principal_discriminator": schema.StringAttribute{
				Required:      true,
				PlanModifiers: requiresReplace,
				Description:   "Principal type: User, Group, or ServicePrincipal.",
			},
			"role_name": schema.StringAttribute{
				Required:    true,
				Description: "Integration role: Owner, Contributor, Reader, or IdentityAccessManager.",
			},
		},
	}
}

func (r *integrationRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data integrationRoleAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}
	result, httpError := r.client.Post(roleAssignmentEndpoint, jsonMap)
	if httpError != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Error calling POST: "+httpError.Error.Error())
		return
	}
	if err = utils.JsonToPlan(result, &data); err != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data integrationRoleAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", roleAssignmentEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}
	if err := utils.JsonToPlan(result, &data); err != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data integrationRoleAssignmentModel
	var state integrationRoleAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}
	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", roleAssignmentEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Error calling PUT: "+httpError.Error.Error())
		return
	}
	if err = utils.JsonToPlan(result, &data); err != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data integrationRoleAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", roleAssignmentEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Error calling DELETE: "+httpError.Error.Error())
		return
	}
}

func (r *integrationRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data integrationRoleAssignmentModel
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", roleAssignmentEndpoint, req.ID))
	if httpError != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}
	if err := utils.JsonToPlan(result, &data); err != nil {
		resp.Diagnostics.AddError(roleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
