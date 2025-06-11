package module_param

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

var moduleParamFromSecretScopedToModuleDefaultError = fmt.Sprintf("snapcd_module_param_from_secret_scoped_to_module error")

var moduleParamFromSecretScopedToModuleEndpoint = "/api/ModuleParamFromSecretScopedToModule"

var _ resource.Resource = (*moduleParamFromSecretScopedToModuleResource)(nil)

func ModuleParamFromSecretScopedToModuleResource() resource.Resource {
	return &moduleParamFromSecretScopedToModuleResource{}
}

type moduleParamFromSecretScopedToModuleResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *moduleParamFromSecretScopedToModuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *moduleParamFromSecretScopedToModuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_param_from_secret_scoped_to_module"
}

type moduleParamFromSecretScopedToModuleModel struct {
	Name                    types.String `tfsdk:"name"`
	Id                      types.String `tfsdk:"id"`
	ModuleId                types.String `tfsdk:"module_id"`
	Type                    types.String `tfsdk:"type"`
	SecretScopedToModuleId  types.String `tfsdk:"secret_scoped_to_module_id"`
}

func (r *moduleParamFromSecretScopedToModuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Module Inputs (Parameters) --- Manages a Module Param (From Secret Scoped To Module) in Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescSharedId + "Module Param (From Secret Scoped To Module).",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescSharedName1 + "Module Param (From Secret Scoped To Module). " + DescSharedName2,
			},
			"secret_scoped_to_module_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the Secret Scoped To Module to take as input.",
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("String", "NotString")},
				Default:     stringdefault.StaticString("String"),
				Description: DescSharedSecretType,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescSharedModuleId1 + "Module Param (From Secret Scoped To Module)" + DescSharedModuleId2,
			},
		},
	}
}

func (r *moduleParamFromSecretScopedToModuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data moduleParamFromSecretScopedToModuleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(moduleParamFromSecretScopedToModuleEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleParamFromSecretScopedToModuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data moduleParamFromSecretScopedToModuleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", moduleParamFromSecretScopedToModuleEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleParamFromSecretScopedToModuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data moduleParamFromSecretScopedToModuleModel
	var state moduleParamFromSecretScopedToModuleModel

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
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", moduleParamFromSecretScopedToModuleEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleParamFromSecretScopedToModuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data moduleParamFromSecretScopedToModuleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", moduleParamFromSecretScopedToModuleEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *moduleParamFromSecretScopedToModuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data moduleParamFromSecretScopedToModuleModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", moduleParamFromSecretScopedToModuleEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleParamFromSecretScopedToModuleDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}