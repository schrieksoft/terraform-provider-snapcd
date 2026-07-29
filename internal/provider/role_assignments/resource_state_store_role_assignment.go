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

var stateStoreRoleAssignmentDefaultError = "snapcd_state_store_role_assignment error"

var stateStoreRoleAssignmentEndpoint = "/StateStoreRoleAssignment"

var _ resource.Resource = (*stateStoreRoleAssignmentResource)(nil)

func StateStoreRoleAssignmentResource() resource.Resource {
	return &stateStoreRoleAssignmentResource{}
}

type stateStoreRoleAssignmentResource struct {
	client *snapcd.Client
}

func (r *stateStoreRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *stateStoreRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_state_store_role_assignment"
}

type stateStoreRoleAssignmentModel struct {
	Id                     types.String `tfsdk:"id"`
	StateStoreId           types.String `tfsdk:"state_store_id"`
	PrincipalId            types.String `tfsdk:"principal_id"`
	PrincipalDiscriminator types.String `tfsdk:"principal_discriminator"`
	RoleName               types.String `tfsdk:"role_name"`
}

func (r *stateStoreRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Role Assignments --- Manages a State Store Role Assignment in Snap CD.` + "\n\n## Required permissions\n\n" + openapidocs.ResourcePermissions["StateStoreRoleAssignment"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: openapidocs.StateStoreRoleAssignmentUpdateDto_Id,
			},
			"state_store_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.StateStoreRoleAssignmentUpdateDto_StateStoreId,
			},
			"principal_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.StateStoreRoleAssignmentUpdateDto_PrincipalId,
			},
			"principal_discriminator": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					// Intentionally narrower than the spec enum: `Base` is an internal discriminator.
					stringvalidator.OneOf("User", "ServicePrincipal", "Group"),
				},
				Description: openapidocs.StateStoreRoleAssignmentUpdateDto_PrincipalDiscriminator,
			},
			"role_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.StateStoreRoleValues...),
				},
				Description: openapidocs.StateStoreRoleAssignmentUpdateDto_RoleName,
			},
		},
	}
}

func (r *stateStoreRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data stateStoreRoleAssignmentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(stateStoreRoleAssignmentEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})

	if err != nil {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stateStoreRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data stateStoreRoleAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", stateStoreRoleAssignmentEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stateStoreRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data stateStoreRoleAssignmentModel
	var state stateStoreRoleAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", stateStoreRoleAssignmentEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stateStoreRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data stateStoreRoleAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", stateStoreRoleAssignmentEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *stateStoreRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data stateStoreRoleAssignmentModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", stateStoreRoleAssignmentEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(stateStoreRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
