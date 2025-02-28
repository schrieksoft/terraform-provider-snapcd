// Copyright (c) HashiCorp, Inc.

package provider

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

var servicePrincipalDefaultError = fmt.Sprintf("snapcd_service_principal error")

var servicePrincipalEndpoint = "/api/Identity/ServicePrincipal"


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

type servicePrincipalModel struct {
	Id               types.String `tfsdk:"id"`
	ClientId         types.String `tfsdk:"client_id"`
	ClientSecret     types.String `tfsdk:"client_secret"`
}

func (r *servicePrincipalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				Required: true,
			},
			"client_secret": schema.StringAttribute{
				Required: true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
		},
	}
}

func (r *servicePrincipalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data servicePrincipalModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, err := r.client.Post(servicePrincipalEndpoint, jsonMap)
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})

	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *servicePrincipalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data servicePrincipalModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := fmt.Sprintf("%s/WithVerifySecret/%s", servicePrincipalEndpoint, data.Id.ValueString())

	// Append the query parameter only if ClientSecret is not empty
	if data.ClientSecret.ValueString() != "" {
		endpoint = fmt.Sprintf("%s?secret=%s", endpoint, data.ClientSecret.ValueString())
	}

	result, err := r.client.Get(endpoint)
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}


func (r *servicePrincipalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data servicePrincipalModel
	var state servicePrincipalModel

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
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, err := r.client.Put(fmt.Sprintf("%s/%s", servicePrincipalEndpoint, state.Id.ValueString()), jsonMap)
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *servicePrincipalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data servicePrincipalModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.client.Delete(fmt.Sprintf("%s/%s", servicePrincipalEndpoint, data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *servicePrincipalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data servicePrincipalModel

	result, err := r.client.Get(fmt.Sprintf("%s/%s", servicePrincipalEndpoint, req.ID))
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data, []string{"client_secret"})
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	data.ClientSecret = types.StringValue("")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
