package namespace_env_var

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

var namespaceEnvVarFromLiteralDefaultError = fmt.Sprintf("snapcd_namespace_env_var_from_literal error")

var namespaceEnvVarFromLiteralEndpoint = "/api/NamespaceEnvVarFromLiteral"

var _ resource.Resource = (*namespaceEnvVarFromLiteralResource)(nil)

func NamespaceEnvVarFromLiteralResource() resource.Resource {
	return &namespaceEnvVarFromLiteralResource{}
}

type namespaceEnvVarFromLiteralResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *namespaceEnvVarFromLiteralResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *namespaceEnvVarFromLiteralResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_env_var_from_literal"
}

type namespaceEnvVarFromLiteralModel struct {
	Name         types.String `tfsdk:"name"`
	Id           types.String `tfsdk:"id"`
	LiteralValue types.String `tfsdk:"literal_value"`
	Type         types.String `tfsdk:"type"`
	UsageMode    types.String `tfsdk:"usage_mode"`
	NamespaceId  types.String `tfsdk:"namespace_id"`
}

func (r *namespaceEnvVarFromLiteralResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Namespace Inputs (Env Vars) --- Manages a Namespace Env Var (From Literal) in Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescSharedId + "Namespace Env Var (From Literal).",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescSharedName1 + "Namespace Env Var (From Literal). " + DescSharedName2,
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("String", "NotString")},
				Default:     stringdefault.StaticString("String"),
				Description: DescSharedLiteralType,
			},
			"literal_value": schema.StringAttribute{
				Required:    true,
				Description: DescSharedLiteralValue,
			},
			"usage_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("UseIfSelected", "UseByDefault"),
				},
				Default: stringdefault.StaticString("UseByDefault"),
				Description: DescSharedUsage,		
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedNamespaceId1 + "Namespace Param (From Literal)" + DescSharedNamespaceId2,
			},
		},
	}
}

func (r *namespaceEnvVarFromLiteralResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data namespaceEnvVarFromLiteralModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(namespaceEnvVarFromLiteralEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceEnvVarFromLiteralResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data namespaceEnvVarFromLiteralModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceEnvVarFromLiteralEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceEnvVarFromLiteralResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data namespaceEnvVarFromLiteralModel
	var state namespaceEnvVarFromLiteralModel

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
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", namespaceEnvVarFromLiteralEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceEnvVarFromLiteralResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data namespaceEnvVarFromLiteralModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", namespaceEnvVarFromLiteralEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *namespaceEnvVarFromLiteralResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data namespaceEnvVarFromLiteralModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceEnvVarFromLiteralEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceEnvVarFromLiteralDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
