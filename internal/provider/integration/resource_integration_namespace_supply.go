package integration

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

var integrationNamespaceSupplyDefaultError = fmt.Sprintf("snapcd_integration_namespace_supply error")

var integrationNamespaceSupplyEndpoint = "/IntegrationNamespaceSupply"

var _ resource.Resource = (*integrationNamespaceSupplyResource)(nil)

func IntegrationNamespaceSupplyResource() resource.Resource {
	return &integrationNamespaceSupplyResource{}
}

type integrationNamespaceSupplyResource struct {
	client *snapcd.Client
}

func (r *integrationNamespaceSupplyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *integrationNamespaceSupplyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_namespace_supply"
}

// ! Category: Integration
type integrationNamespaceSupplyModel struct {
	Id            types.String `tfsdk:"id"`
	NamespaceId   types.String `tfsdk:"namespace_id"`
	IntegrationId types.String `tfsdk:"integration_id"`
}

func (r *integrationNamespaceSupplyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Integrations --- Manages an Integration Namespace Supply in Snap CD. Supplies the integration to every module in the given namespace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of the Integration Namespace Supply.",
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the Namespace the integration is supplied to.",
			},
			"integration_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the Integration that is supplied to the Namespace.",
			},
		},
	}
}

func (r *integrationNamespaceSupplyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data integrationNamespaceSupplyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(integrationNamespaceSupplyEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Error calling POST, unexpected error: "+httpError.Error.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationNamespaceSupplyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data integrationNamespaceSupplyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", integrationNamespaceSupplyEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Error calling GET, unexpected error: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationNamespaceSupplyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data integrationNamespaceSupplyModel
	var state integrationNamespaceSupplyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", integrationNamespaceSupplyEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Error calling PUT, unexpected error: "+httpError.Error.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationNamespaceSupplyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data integrationNamespaceSupplyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", integrationNamespaceSupplyEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpError != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Error calling DELETE, unexpected error: "+httpError.Error.Error())
		return
	}
}

func (r *integrationNamespaceSupplyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data integrationNamespaceSupplyModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", integrationNamespaceSupplyEndpoint, req.ID))
	if httpError != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Error calling GET, unexpected error: "+httpError.Error.Error())
		return
	}

	err := utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(integrationNamespaceSupplyDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
