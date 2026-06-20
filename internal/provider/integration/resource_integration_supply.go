package integration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
)

var supplyDefaultError = "snapcd_integration_supply error"

var _ resource.Resource = (*integrationSupplyResource)(nil)

func IntegrationSupplyResource() resource.Resource {
	return &integrationSupplyResource{}
}

type integrationSupplyResource struct {
	client *snapcd.Client
}

func (r *integrationSupplyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *integrationSupplyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_supply"
}

// ! Category: Integration
type integrationSupplyModel struct {
	Id            types.String `tfsdk:"id"`
	IntegrationId types.String `tfsdk:"integration_id"`
	Scope         types.String `tfsdk:"scope"`
	ScopeId       types.String `tfsdk:"scope_id"`
}

func (r *integrationSupplyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	requiresReplace := []planmodifier.String{stringplanmodifier.RequiresReplace()}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Integrations --- Supplies an integration to a Stack, Namespace, or Module scope.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "Unique ID of the supply.",
			},
			"integration_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: requiresReplace,
				Description:   "ID of the integration.",
			},
			"scope": schema.StringAttribute{
				Required:      true,
				PlanModifiers: requiresReplace,
				Description:   "Scope: Stack, Namespace, or Module.",
			},
			"scope_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: requiresReplace,
				Description:   "ID of the stack/namespace/module the integration is supplied to.",
			},
		},
	}
}

func (r *integrationSupplyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data integrationSupplyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := map[string]interface{}{"Scope": data.Scope.ValueString(), "ScopeId": data.ScopeId.ValueString()}
	result, httpError := r.client.Post(fmt.Sprintf("%s/%s/supplies", integrationEndpoint, data.IntegrationId.ValueString()), body)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(supplyDefaultError, "This supply already exists. Import it to manage with terraform.")
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(supplyDefaultError, "Error calling POST: "+httpError.Error.Error())
		return
	}
	data.Id = types.StringValue(getStr(result, "id"))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationSupplyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data integrationSupplyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s/supplies/%s/%s", integrationEndpoint, data.IntegrationId.ValueString(), data.Scope.ValueString(), data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(supplyDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}
	data.Id = types.StringValue(getStr(result, "id"))
	if sid := getStr(result, "scopeId"); sid != "" {
		data.ScopeId = types.StringValue(sid)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is unreachable — every attribute forces replacement — but the interface requires it.
func (r *integrationSupplyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data integrationSupplyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationSupplyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data integrationSupplyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s/supplies/%s/%s", integrationEndpoint, data.IntegrationId.ValueString(), data.Scope.ValueString(), data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(supplyDefaultError, "Error calling DELETE: "+httpError.Error.Error())
		return
	}
}

// ImportState ID format: "<integration_id>,<scope>,<supply_id>".
func (r *integrationSupplyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(supplyDefaultError, "Import ID must be '<integration_id>,<scope>,<supply_id>'.")
		return
	}
	data := integrationSupplyModel{
		IntegrationId: types.StringValue(parts[0]),
		Scope:         types.StringValue(parts[1]),
		Id:            types.StringValue(parts[2]),
	}
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s/supplies/%s/%s", integrationEndpoint, parts[0], parts[1], parts[2]))
	if httpError != nil {
		resp.Diagnostics.AddError(supplyDefaultError, "Error calling GET: "+httpError.Error.Error())
		return
	}
	data.ScopeId = types.StringValue(getStr(result, "scopeId"))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
