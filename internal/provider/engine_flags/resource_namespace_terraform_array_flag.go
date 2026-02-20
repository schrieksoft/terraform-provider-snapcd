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

var namespaceTerraformArrayFlagDefaultError = fmt.Sprintf("snapcd_namespace_terraform_array_flag error")

var namespaceTerraformArrayFlagEndpoint = "/NamespaceTerraformArrayFlag"

var _ resource.Resource = (*namespaceTerraformArrayFlagResource)(nil)

func NamespaceTerraformArrayFlagResource() resource.Resource {
	return &namespaceTerraformArrayFlagResource{}
}

type namespaceTerraformArrayFlagResource struct {
	client *snapcd.Client
}

func (r *namespaceTerraformArrayFlagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *namespaceTerraformArrayFlagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_terraform_array_flag"
}

type namespaceTerraformArrayFlagModel struct {
	Id          types.String `tfsdk:"id"`
	NamespaceId types.String `tfsdk:"namespace_id"`
	Task        types.String `tfsdk:"task"`
	Flag        types.String `tfsdk:"flag"`
	Value       types.String `tfsdk:"value"`
}

const (
	DescNamespaceTerraformArrayFlagId          = "Unique ID of the Namespace Terraform Array Flag."
	DescNamespaceTerraformArrayFlagNamespaceId = "ID of the parent Namespace."
	DescNamespaceTerraformArrayFlagTask        = "The command task this flag applies to. Valid values: `Init`, `Plan`, `Apply`, `Destroy`, `Output`."
	DescNamespaceTerraformArrayFlagFlag        = "The Terraform CLI array flag name."
	DescNamespaceTerraformArrayFlagValue       = "The value for the flag."
)

func (r *namespaceTerraformArrayFlagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Engine Flags --- Manages a Namespace Terraform Array Flag in Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescNamespaceTerraformArrayFlagId,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceTerraformArrayFlagNamespaceId,
			},
			"task": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceTerraformArrayFlagTask,
				Validators: []validator.String{
					stringvalidator.OneOf(terraformCommandTaskValues...),
				},
			},
			"flag": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceTerraformArrayFlagFlag,
				Validators: []validator.String{
					stringvalidator.OneOf(terraformArrayFlagValues...),
				},
			},
			"value": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceTerraformArrayFlagValue,
			},
		},
	}
}

func (r *namespaceTerraformArrayFlagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data namespaceTerraformArrayFlagModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(namespaceTerraformArrayFlagEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceTerraformArrayFlagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data namespaceTerraformArrayFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceTerraformArrayFlagEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceTerraformArrayFlagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data namespaceTerraformArrayFlagModel
	var state namespaceTerraformArrayFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", namespaceTerraformArrayFlagEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceTerraformArrayFlagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data namespaceTerraformArrayFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", namespaceTerraformArrayFlagEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *namespaceTerraformArrayFlagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data namespaceTerraformArrayFlagModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", namespaceTerraformArrayFlagEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(namespaceTerraformArrayFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
