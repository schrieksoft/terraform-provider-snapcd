package identity

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"
	"terraform-provider-snapcd/internal/validators"

	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var servicePrincipalResourceDefaultError = "snapcd_service_principal error"

var _ resource.Resource = (*servicePrincipalResource)(nil)

func ServicePrincipalResource() resource.Resource {
	return &servicePrincipalResource{}
}

type servicePrincipalResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *servicePrincipalResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *servicePrincipalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_principal"
}

type servicePrincipalResourceModel struct {
	Id          types.String `tfsdk:"id"`
	ClientId    types.String `tfsdk:"client_id"`
	DisplayName types.String `tfsdk:"display_name"`
	IsDisabled  types.Bool   `tfsdk:"is_disabled"`
	Scopes      types.List   `tfsdk:"scopes"`
}

func (r *servicePrincipalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Identity Access Management --- Manages a Service Principal in Snap CD.

The Service Principal is created without a client secret, and this resource cannot set one: a secret managed through Terraform would be written to state by whatever generates it and whatever consumes it. Give it a ` + "`snapcd_service_principal_federated_credential`" + ` instead, so the workload authenticates with a token its own platform issues and no secret exists.

A secret can still be set in the Dashboard for a Service Principal that is not managed here.` + "\n\n## Required permissions\n\n" + openapidocs.ResourcePermissions["ServicePrincipal"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: openapidocs.ServicePrincipalReadDto_Id,
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ServicePrincipalCreateDto_ClientId,
				Validators: []validator.String{
					validators.NotBlank(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: openapidocs.ServicePrincipalMetadataReadDto_DisplayName,
			},
			"is_disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.ServicePrincipalCreateDto_IsDisabled,
				Default:     booldefault.StaticBool(false),
			},
			"scopes": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: openapidocs.ServicePrincipalCreateDto_Scopes,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("snapcd_scope")})),
			},
		},
	}
}

func (r *servicePrincipalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data servicePrincipalResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id", "displayname"})
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(servicePrincipalEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *servicePrincipalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data servicePrincipalResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Metadata, not the full record: it carries the client ID and display name without the
	// stored secret hash, and needs only the metadata reader permission.
	result, httpError := r.client.Get(fmt.Sprintf("%s/Metadata/%s", servicePrincipalEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		utils.ExternallyDeleted(ctx, &resp.Diagnostics, "snapcd_service_principal", data.Id.ValueString())
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
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *servicePrincipalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data servicePrincipalResourceModel
	var state servicePrincipalResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id", "displayname"})
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", servicePrincipalEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *servicePrincipalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data servicePrincipalResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", servicePrincipalEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
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
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *servicePrincipalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data servicePrincipalResourceModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/Metadata/%s", servicePrincipalEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalResourceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
