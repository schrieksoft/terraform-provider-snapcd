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

var namespaceRoleAssignmentDefaultError = fmt.Sprintf("snapcd_namespace_role_assignment error")

var namespaceRoleAssignmentEndpoint = "/NamespaceRoleAssignment"

var _ resource.Resource = (*namespaceRoleAssignmentResource)(nil)

func NamespaceRoleAssignmentResource() resource.Resource {
	return &namespaceRoleAssignmentResource{}
}

type namespaceRoleAssignmentResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *namespaceRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *namespaceRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_role_assignment"
}

type namespaceRoleAssignmentModel struct {
	Id                     types.String `tfsdk:"id"`
	NamespaceId            types.String `tfsdk:"namespace_id"`
	PrincipalId            types.String `tfsdk:"principal_id"`
	PrincipalDiscriminator types.String `tfsdk:"principal_discriminator"`
	RoleName               types.String `tfsdk:"role_name"`
}

func (r *namespaceRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Role Assignments --- Manages a Namespace Role Assignment in Snap CD.` + "\n\n## Required permissions\n\n" + openapidocs.ResourcePermissions["NamespaceRoleAssignment"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: openapidocs.NamespaceRoleAssignmentReadDto_Id,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceRoleAssignmentReadDto_NamespaceId,
			},
			"principal_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.NamespaceRoleAssignmentReadDto_PrincipalId,
			},
			"principal_discriminator": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					// Intentionally narrower than the spec enum: `Base` is an internal discriminator.
					stringvalidator.OneOf("User", "ServicePrincipal", "Group"),
				},
				Description: openapidocs.NamespaceRoleAssignmentReadDto_PrincipalDiscriminator,
			},
			"role_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(openapidocs.NamespaceRoleValues...),
				},
				Description: openapidocs.NamespaceRoleAssignmentReadDto_RoleName,
			},
		},
	}
}

func (r *namespaceRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data namespaceRoleAssignmentModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(namespaceRoleAssignmentEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})

	if err != nil {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data namespaceRoleAssignmentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceRoleAssignmentEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data namespaceRoleAssignmentModel
	var state namespaceRoleAssignmentModel

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
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", namespaceRoleAssignmentEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data namespaceRoleAssignmentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", namespaceRoleAssignmentEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *namespaceRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data namespaceRoleAssignmentModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceRoleAssignmentEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(namespaceRoleAssignmentDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
