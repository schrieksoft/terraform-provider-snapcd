package integration_events

import (
	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var moduleIntegrationEventDefaultError = fmt.Sprintf("snapcd_module_integration_event error")

var _ resource.Resource = (*moduleIntegrationEventResource)(nil)

func ModuleIntegrationEventResource() resource.Resource {
	return &moduleIntegrationEventResource{}
}

type moduleIntegrationEventResource struct {
	client *snapcd.Client
}

func (r *moduleIntegrationEventResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *moduleIntegrationEventResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_integration_event"
}

// ! Category: Integration
type moduleIntegrationEventModel struct {
	Id            types.String `tfsdk:"id"`
	ModuleId      types.String `tfsdk:"module_id"`
	IntegrationId types.String `tfsdk:"integration_id"`
	Trigger       types.String `tfsdk:"trigger"`
	Template      types.String `tfsdk:"template"`
	Filter        types.String `tfsdk:"filter"`
	IsDisabled    types.Bool   `tfsdk:"is_disabled"`
}

func (r *moduleIntegrationEventResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Integration Events --- Manages a module-scoped Integration Event in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescEventId,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescEventModuleId,
			},
			"integration_id": schema.StringAttribute{
				Required:    true,
				Description: DescEventIntegrationId,
			},
			"trigger": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(integrationTriggerValues...),
				},
				Description: DescEventTrigger,
			},
			"template": schema.StringAttribute{
				Optional:    true,
				Description: DescEventTemplate,
			},
			"filter": schema.StringAttribute{
				Optional:    true,
				Description: DescEventFilter,
			},
			"is_disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: DescEventIsDisabled,
			},
		},
	}
}

func (r *moduleIntegrationEventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data moduleIntegrationEventModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(moduleIntegrationEventEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Error calling POST, unexpected error: "+httpError.Error.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleIntegrationEventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data moduleIntegrationEventModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", moduleIntegrationEventEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Error calling GET, unexpected error: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleIntegrationEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data moduleIntegrationEventModel
	var state moduleIntegrationEventModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", moduleIntegrationEventEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Error calling PUT, unexpected error: "+httpError.Error.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *moduleIntegrationEventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data moduleIntegrationEventModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", moduleIntegrationEventEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Error calling DELETE, unexpected error: "+httpError.Error.Error())
		return
	}
}

func (r *moduleIntegrationEventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data moduleIntegrationEventModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", moduleIntegrationEventEndpoint, req.ID))
	if httpError != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Error calling GET, unexpected error: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(moduleIntegrationEventDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
