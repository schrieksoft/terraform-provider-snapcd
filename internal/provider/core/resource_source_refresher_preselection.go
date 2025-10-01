package core

import (
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

var sourceRefresherPreselectionDefaultError = fmt.Sprintf("snapcd_sourceRefresherPreselection error")

var sourceRefresherPreselectionEndpoint = "/SourceRefresherPreselection"

var _ resource.Resource = (*sourceRefresherPreselectionResource)(nil)

func SourceRefresherPreselectionResource() resource.Resource {
	return &sourceRefresherPreselectionResource{}
}

type sourceRefresherPreselectionResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *sourceRefresherPreselectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *sourceRefresherPreselectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_refresher_preselection"
}

// ! Category: Source Refresher Preselection
type sourceRefresherPreselectionModel struct {
	Id                     types.String `tfsdk:"id"`
	RunnerPoolId           types.String `tfsdk:"runner_pool_id"`
	RunnerSelfDeclaredName types.String `tfsdk:"runner_self_declared_name"`
	SourceUrl              types.String `tfsdk:"source_url"`
}

const (
	DescSourceRefresherPreselectionId                         = "Unique ID of the Source Refresher Preselection."
	DescSourceRefresherPreselectionSourceUrl                  = "Unique Source URL to which a Runner Pool (or specific Runner within the Runner Pool based on `runner_self_declared_name`) is assigned as the preselected 'refresher'."
	DescSourceRefresherPreselectionRunnerPoolId               = "ID of the Runner Pool to preselect as 'refresher' for the given Source URL. Messages requesting a source refresh will be sent to this Runner Pool's service bus endpoint."
	DescSourceRefresherPreselectionRunnerPoolSelfDeclaredName = "Self-declared name of a Runner (within the Runner Pool defined by `runner_pool_id`) to preselect as 'refresher'. If this is set, then messages will be sent to this Runner's specifc service bus endpoint, instead of sending to the Runner Pool's shared endpoint."
)

func (r *sourceRefresherPreselectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Source Refresher Preselections --- Manages a Source Refresher Preselection in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescSourceRefresherPreselectionId,
			},
			"source_url": schema.StringAttribute{
				Required:    true,
				Description: DescSourceRefresherPreselectionSourceUrl,
			},
			"runner_pool_id": schema.StringAttribute{
				Required:    true,
				Description: DescSourceRefresherPreselectionRunnerPoolId,
			},
			"runner_self_declared_name": schema.StringAttribute{
				Optional:    true,
				Description: DescSourceRefresherPreselectionRunnerPoolSelfDeclaredName,
			},
		},
	}
}

func (r *sourceRefresherPreselectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data sourceRefresherPreselectionModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(sourceRefresherPreselectionEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceRefresherPreselectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data sourceRefresherPreselectionModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", sourceRefresherPreselectionEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceRefresherPreselectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data sourceRefresherPreselectionModel
	var state sourceRefresherPreselectionModel

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
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", sourceRefresherPreselectionEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceRefresherPreselectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data sourceRefresherPreselectionModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", sourceRefresherPreselectionEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *sourceRefresherPreselectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data sourceRefresherPreselectionModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", sourceRefresherPreselectionEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(sourceRefresherPreselectionDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
