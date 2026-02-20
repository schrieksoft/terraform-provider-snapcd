package engine_flags

import (
	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var namespaceTerraformFlagDefaultError = fmt.Sprintf("snapcd_namespace_terraform_flag error")

var namespaceTerraformFlagEndpoint = "/NamespaceTerraformFlag"

var _ resource.Resource = (*namespaceTerraformFlagResource)(nil)

func NamespaceTerraformFlagResource() resource.Resource {
	return &namespaceTerraformFlagResource{}
}

type namespaceTerraformFlagResource struct {
	client *snapcd.Client
}

func (r *namespaceTerraformFlagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *namespaceTerraformFlagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_terraform_flag"
}

type namespaceTerraformFlagModel struct {
	Id          types.String `tfsdk:"id"`
	NamespaceId types.String `tfsdk:"namespace_id"`
	Task        types.String `tfsdk:"task"`
	Flag        types.String `tfsdk:"flag"`
	Value       types.String `tfsdk:"value"`
}

const (
	DescNamespaceTerraformFlagId          = "Unique ID of the Namespace Terraform Flag."
	DescNamespaceTerraformFlagNamespaceId = "ID of the parent Namespace."
	DescNamespaceTerraformFlagTask        = "The command task this flag applies to. Valid values: `Init`, `Plan`, `Apply`, `Destroy`, `Output`."
	DescNamespaceTerraformFlagFlag        = "The Terraform CLI flag name."
	DescNamespaceTerraformFlagValue       = "The value for the flag. Optional for boolean flags."
)

func (r *namespaceTerraformFlagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Engine Flags --- Manages a Namespace Terraform Flag in Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescNamespaceTerraformFlagId,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceTerraformFlagNamespaceId,
			},
			"task": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceTerraformFlagTask,
				Validators: []validator.String{
					stringvalidator.OneOf(terraformCommandTaskValues...),
				},
			},
			"flag": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceTerraformFlagFlag,
				Validators: []validator.String{
					stringvalidator.OneOf(terraformFlagValues...),
				},
			},
			"value": schema.StringAttribute{
				Optional:    true,
				Description: DescNamespaceTerraformFlagValue,
			},
		},
	}
}

func (r *namespaceTerraformFlagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data namespaceTerraformFlagModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(namespaceTerraformFlagEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceTerraformFlagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data namespaceTerraformFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceTerraformFlagEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceTerraformFlagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data namespaceTerraformFlagModel
	var state namespaceTerraformFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", namespaceTerraformFlagEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceTerraformFlagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data namespaceTerraformFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", namespaceTerraformFlagEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *namespaceTerraformFlagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data namespaceTerraformFlagModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceTerraformFlagEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
