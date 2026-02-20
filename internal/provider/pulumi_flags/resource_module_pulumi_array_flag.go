package pulumi_flags

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

var modulePulumiArrayFlagDefaultError = fmt.Sprintf("snapcd_module_pulumi_array_flag error")

var modulePulumiArrayFlagEndpoint = "/ModulePulumiArrayFlag"

var _ resource.Resource = (*modulePulumiArrayFlagResource)(nil)

func ModulePulumiArrayFlagResource() resource.Resource {
	return &modulePulumiArrayFlagResource{}
}

type modulePulumiArrayFlagResource struct {
	client *snapcd.Client
}

func (r *modulePulumiArrayFlagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *modulePulumiArrayFlagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_pulumi_array_flag"
}

type modulePulumiArrayFlagModel struct {
	Id       types.String `tfsdk:"id"`
	ModuleId types.String `tfsdk:"module_id"`
	Task     types.String `tfsdk:"task"`
	Flag     types.String `tfsdk:"flag"`
	Value    types.String `tfsdk:"value"`
}

const (
	DescModulePulumiArrayFlagId       = "Unique ID of the Module Pulumi Array Flag."
	DescModulePulumiArrayFlagModuleId = "ID of the parent Module."
	DescModulePulumiArrayFlagTask     = "The command task this flag applies to. Valid values: `Init`, `Plan`, `Apply`, `Destroy`, `Output`."
	DescModulePulumiArrayFlagFlag     = "The Pulumi CLI array flag name."
	DescModulePulumiArrayFlagValue    = "The value for the flag."
)

func (r *modulePulumiArrayFlagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Pulumi Flags --- Manages a Module Pulumi Array Flag in Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescModulePulumiArrayFlagId,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescModulePulumiArrayFlagModuleId,
			},
			"task": schema.StringAttribute{
				Required:    true,
				Description: DescModulePulumiArrayFlagTask,
				Validators: []validator.String{
					stringvalidator.OneOf(pulumiCommandTaskValues...),
				},
			},
			"flag": schema.StringAttribute{
				Required:    true,
				Description: DescModulePulumiArrayFlagFlag,
				Validators: []validator.String{
					stringvalidator.OneOf(pulumiArrayFlagValues...),
				},
			},
			"value": schema.StringAttribute{
				Required:    true,
				Description: DescModulePulumiArrayFlagValue,
			},
		},
	}
}

func (r *modulePulumiArrayFlagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data modulePulumiArrayFlagModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(modulePulumiArrayFlagEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *modulePulumiArrayFlagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data modulePulumiArrayFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", modulePulumiArrayFlagEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *modulePulumiArrayFlagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data modulePulumiArrayFlagModel
	var state modulePulumiArrayFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", modulePulumiArrayFlagEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *modulePulumiArrayFlagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data modulePulumiArrayFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", modulePulumiArrayFlagEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *modulePulumiArrayFlagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data modulePulumiArrayFlagModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", modulePulumiArrayFlagEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiArrayFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
