package integration

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

var stackIntegrationEventDefaultError = fmt.Sprintf("snapcd_stack_integration_event error")

var _ resource.Resource = (*stackIntegrationEventResource)(nil)

func StackIntegrationEventResource() resource.Resource {
	return &stackIntegrationEventResource{}
}

type stackIntegrationEventResource struct {
	client *snapcd.Client
}

func (r *stackIntegrationEventResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *stackIntegrationEventResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack_integration_event"
}

// ! Category: Integration
type stackIntegrationEventModel struct {
	Id            types.String `tfsdk:"id"`
	StackId       types.String `tfsdk:"stack_id"`
	IntegrationId types.String `tfsdk:"integration_id"`
	Trigger       types.String `tfsdk:"trigger"`
	Template      types.String `tfsdk:"template"`
	Filter        types.String `tfsdk:"filter"`
	IsDisabled    types.Bool   `tfsdk:"is_disabled"`
}

func (r *stackIntegrationEventResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Integrations --- Manages a stack-scoped Integration Event in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescEventId,
			},
			"stack_id": schema.StringAttribute{
				Required:    true,
				Description: DescEventStackId,
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

func (r *stackIntegrationEventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data stackIntegrationEventModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(stackIntegrationEventEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Error calling POST, unexpected error: "+httpError.Error.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stackIntegrationEventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data stackIntegrationEventModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", stackIntegrationEventEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Error calling GET, unexpected error: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stackIntegrationEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data stackIntegrationEventModel
	var state stackIntegrationEventModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", stackIntegrationEventEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Error calling PUT, unexpected error: "+httpError.Error.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stackIntegrationEventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data stackIntegrationEventModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", stackIntegrationEventEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Error calling DELETE, unexpected error: "+httpError.Error.Error())
		return
	}
}

func (r *stackIntegrationEventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data stackIntegrationEventModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", stackIntegrationEventEndpoint, req.ID))
	if httpError != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Error calling GET, unexpected error: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(stackIntegrationEventDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
