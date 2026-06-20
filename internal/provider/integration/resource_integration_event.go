package integration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
)

var eventDefaultError = "snapcd_integration_event error"
var eventEndpoint = "/IntegrationEvent"

var _ resource.Resource = (*integrationEventResource)(nil)

func IntegrationEventResource() resource.Resource {
	return &integrationEventResource{}
}

type integrationEventResource struct {
	client *snapcd.Client
}

func (r *integrationEventResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *integrationEventResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_event"
}

// ! Category: Integration
type integrationEventModel struct {
	Id            types.String `tfsdk:"id"`
	Scope         types.String `tfsdk:"scope"`
	ScopeId       types.String `tfsdk:"scope_id"`
	IntegrationId types.String `tfsdk:"integration_id"`
	Trigger       types.String `tfsdk:"trigger"`
	Template      types.String `tfsdk:"template"`
	Filter        types.String `tfsdk:"filter"`
	IsDisabled    types.Bool   `tfsdk:"is_disabled"`
}

func (r *integrationEventResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	requiresReplace := []planmodifier.String{stringplanmodifier.RequiresReplace()}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Integrations --- Subscribes a trigger on a scope to an integration (the demand side).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "Unique ID of the event subscription.",
			},
			"scope": schema.StringAttribute{
				Required:      true,
				PlanModifiers: requiresReplace,
				Description:   "Scope: Organization, Stack, Namespace, or Module.",
			},
			"scope_id": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: requiresReplace,
				Description:   "ID of the stack/namespace/module (omit for Organization scope).",
			},
			"integration_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the target integration.",
			},
			"trigger": schema.StringAttribute{
				Required:    true,
				Description: "Trigger this subscription fires on (e.g. JobFailed, MissionMilestoneReported).",
			},
			"template": schema.StringAttribute{
				Optional:    true,
				Description: "Optional message template ({{token}} substitution). Omit to use the default for the trigger.",
			},
			"filter": schema.StringAttribute{
				Optional:    true,
				Description: "Optional filter expression.",
			},
			"is_disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the subscription is disabled.",
			},
		},
	}
}

func optStr(v types.String) interface{} {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return v.ValueString()
}

func (r *integrationEventResource) applyResponse(result map[string]interface{}, data *integrationEventModel) {
	data.Id = types.StringValue(getStr(result, "id"))
	data.Scope = types.StringValue(getStr(result, "scope"))
	data.IntegrationId = types.StringValue(getStr(result, "integrationId"))
	data.Trigger = types.StringValue(getStr(result, "trigger"))
	data.Template = getOptStr(result, "template")
	data.Filter = getOptStr(result, "filter")
	data.IsDisabled = types.BoolValue(getBool(result, "isDisabled"))
	if sid := getOptStr(result, "scopeId"); !sid.IsNull() {
		data.ScopeId = sid
	}
}

func (r *integrationEventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data integrationEventModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := map[string]interface{}{
		"Scope":         data.Scope.ValueString(),
		"ScopeId":       optStr(data.ScopeId),
		"IntegrationId": data.IntegrationId.ValueString(),
		"Trigger":       data.Trigger.ValueString(),
		"Template":      optStr(data.Template),
		"Filter":        optStr(data.Filter),
		"IsDisabled":    data.IsDisabled.ValueBool(),
	}
	result, httpError := r.client.Post(eventEndpoint, body)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(eventDefaultError, "This event subscription already exists. Import it to manage with terraform.")
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(eventDefaultError, "Error calling POST: "+httpError.Error.Error())
		return
	}
	r.applyResponse(result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationEventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data integrationEventModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s/%s", eventEndpoint, data.Scope.ValueString(), data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(eventDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}
	r.applyResponse(result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data integrationEventModel
	var state integrationEventModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := map[string]interface{}{
		"IntegrationId": data.IntegrationId.ValueString(),
		"Trigger":       data.Trigger.ValueString(),
		"Template":      optStr(data.Template),
		"Filter":        optStr(data.Filter),
		"IsDisabled":    data.IsDisabled.ValueBool(),
	}
	_, httpError := r.client.Put(fmt.Sprintf("%s/%s/%s", eventEndpoint, state.Scope.ValueString(), state.Id.ValueString()), body)
	if httpError != nil {
		resp.Diagnostics.AddError(eventDefaultError, "Error calling PUT: "+httpError.Error.Error())
		return
	}
	// PUT returns NoContent; the plan is the desired state.
	data.Id = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationEventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data integrationEventModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s/%s", eventEndpoint, data.Scope.ValueString(), data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(eventDefaultError, "Error calling DELETE: "+httpError.Error.Error())
		return
	}
}

// ImportState ID format: "<scope>,<event_id>".
func (r *integrationEventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(eventDefaultError, "Import ID must be '<scope>,<event_id>'.")
		return
	}
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s/%s", eventEndpoint, parts[0], parts[1]))
	if httpError != nil {
		resp.Diagnostics.AddError(eventDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}
	var data integrationEventModel
	r.applyResponse(result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getOptStr(m map[string]interface{}, key string) types.String {
	if v, ok := lookup(m, key); ok {
		if s, ok := v.(string); ok {
			return types.StringValue(s)
		}
	}
	return types.StringNull()
}
