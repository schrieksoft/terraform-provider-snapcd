package output

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"

	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var outputNamespaceSupplyDefaultError = fmt.Sprintf("snapcd_output_namespace_supply error")

var outputNamespaceSupplyEndpoint = "/OutputNamespaceSupply"

var _ resource.Resource = (*outputNamespaceSupplyResource)(nil)

func OutputNamespaceSupplyResource() resource.Resource {
	return &outputNamespaceSupplyResource{}
}

type outputNamespaceSupplyResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *outputNamespaceSupplyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *outputNamespaceSupplyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_output_namespace_supply"
}

type outputNamespaceSupplyModel struct {
	Id       types.String `tfsdk:"id"`
	ModuleId types.String `tfsdk:"module_id"`
	OutputName types.String `tfsdk:"output_name"`
	ConsumingNamespaceId types.String `tfsdk:"consuming_namespace_id"`
}

func (r *outputNamespaceSupplyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Outputs --- Manages an Output Namespace Supply in Snap CD." + "\n\n## Required permissions\n\n" + openapidocs.ResourcePermissions["OutputNamespaceSupply"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: openapidocs.OutputNamespaceSupplyReadDto_Id,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.OutputNamespaceSupplyCreateDto_ModuleId,
			},
			"output_name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.OutputNamespaceSupplyCreateDto_OutputName,
			},
			"consuming_namespace_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.OutputNamespaceSupplyCreateDto_ConsumingNamespaceId,
			},
		},
	}
}

func (r *outputNamespaceSupplyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data outputNamespaceSupplyModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	if warning := r.client.HoldModule(data.ModuleId.ValueString()); warning != "" {
		resp.Diagnostics.AddWarning(snapcd.ModuleHoldWarningSummary, warning)
	}
	result, httpError := r.client.Post(outputNamespaceSupplyEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *outputNamespaceSupplyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data outputNamespaceSupplyModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", outputNamespaceSupplyEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		utils.ExternallyDeleted(ctx, &resp.Diagnostics, "snapcd_output_namespace_supply", data.Id.ValueString())
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
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *outputNamespaceSupplyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data outputNamespaceSupplyModel
	var state outputNamespaceSupplyModel

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
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	if warning := r.client.HoldModule(data.ModuleId.ValueString()); warning != "" {
		resp.Diagnostics.AddWarning(snapcd.ModuleHoldWarningSummary, warning)
	}
	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", outputNamespaceSupplyEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *outputNamespaceSupplyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data outputNamespaceSupplyModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	if warning := r.client.HoldModule(data.ModuleId.ValueString()); warning != "" {
		resp.Diagnostics.AddWarning(snapcd.ModuleHoldWarningSummary, warning)
	}
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", outputNamespaceSupplyEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *outputNamespaceSupplyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data outputNamespaceSupplyModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", outputNamespaceSupplyEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(outputNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
