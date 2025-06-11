package namespace_param

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

var namespaceParamFromSecretScopedToStackDefaultError = fmt.Sprintf("snapcd_namespace_param_from_secret_scoped_to_stack error")

var namespaceParamFromSecretScopedToStackEndpoint = "/api/NamespaceParamFromSecretScopedToStack"

var _ resource.Resource = (*namespaceParamFromSecretScopedToStackResource)(nil)

func NamespaceParamFromSecretScopedToStackResource() resource.Resource {
	return &namespaceParamFromSecretScopedToStackResource{}
}

type namespaceParamFromSecretScopedToStackResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *namespaceParamFromSecretScopedToStackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *namespaceParamFromSecretScopedToStackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_param_from_secret_scoped_to_stack"
}

type namespaceParamFromSecretScopedToStackModel struct {
	Name                   types.String `tfsdk:"name"`
	Id                     types.String `tfsdk:"id"`
	Type                   types.String `tfsdk:"type"`
	SecretScopedToStackId  types.String `tfsdk:"secret_scoped_to_stack_id"`
	UsageMode              types.String `tfsdk:"usage_mode"`
	NamespaceId            types.String `tfsdk:"namespace_id"`
}

func (r *namespaceParamFromSecretScopedToStackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Namespace Inputs (Parameters) --- Manages a Namespace Param (From Secret Scoped To Stack) in Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescSharedId + "Namespace Param (From Secret Scoped To Stack).",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescSharedName1 + "Namespace Param (From Secret Scoped To Stack). " + DescSharedName2,
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("String", "NotString")},
				Default:     stringdefault.StaticString("String"),
				Description: DescSharedSecretType,
			},
			"secret_scoped_to_stack_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the secret scoped to stack to use for this namespace param.",
			},
			"usage_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("UseIfSelected", "UseByDefault"),
				},
				Default:     stringdefault.StaticString("UseIfSelected"),
				Description: DescSharedUsage,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedNamespaceId1 + "Namespace Param (From Secret Scoped To Stack)" + DescSharedNamespaceId2,
			},
		},
	}
}

func (r *namespaceParamFromSecretScopedToStackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data namespaceParamFromSecretScopedToStackModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(namespaceParamFromSecretScopedToStackEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceParamFromSecretScopedToStackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data namespaceParamFromSecretScopedToStackModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceParamFromSecretScopedToStackEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceParamFromSecretScopedToStackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data namespaceParamFromSecretScopedToStackModel
	var state namespaceParamFromSecretScopedToStackModel

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
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", namespaceParamFromSecretScopedToStackEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceParamFromSecretScopedToStackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data namespaceParamFromSecretScopedToStackModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", namespaceParamFromSecretScopedToStackEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *namespaceParamFromSecretScopedToStackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data namespaceParamFromSecretScopedToStackModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceParamFromSecretScopedToStackEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceParamFromSecretScopedToStackDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}