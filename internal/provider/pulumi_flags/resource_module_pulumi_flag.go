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

var modulePulumiFlagDefaultError = fmt.Sprintf("snapcd_module_pulumi_flag error")

var modulePulumiFlagEndpoint = "/ModulePulumiFlag"

var _ resource.Resource = (*modulePulumiFlagResource)(nil)

func ModulePulumiFlagResource() resource.Resource {
	return &modulePulumiFlagResource{}
}

type modulePulumiFlagResource struct {
	client *snapcd.Client
}

func (r *modulePulumiFlagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *modulePulumiFlagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_pulumi_flag"
}

type modulePulumiFlagModel struct {
	Id       types.String `tfsdk:"id"`
	ModuleId types.String `tfsdk:"module_id"`
	Task     types.String `tfsdk:"task"`
	Flag     types.String `tfsdk:"flag"`
	Value    types.String `tfsdk:"value"`
}

const (
	DescModulePulumiFlagId       = "Unique ID of the Module Pulumi Flag."
	DescModulePulumiFlagModuleId = "ID of the parent Module."
	DescModulePulumiFlagTask     = "The command task this flag applies to. Valid values: `Init`, `Plan`, `Apply`, `Destroy`, `Output`."
	DescModulePulumiFlagFlag     = "The Pulumi CLI flag name."
	DescModulePulumiFlagValue    = "The value for the flag. Optional for boolean flags."
)

func (r *modulePulumiFlagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Pulumi Flags --- Manages a Module Pulumi Flag in Snap CD.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescModulePulumiFlagId,
			},
			"module_id": schema.StringAttribute{
				Required:    true,
				Description: DescModulePulumiFlagModuleId,
			},
			"task": schema.StringAttribute{
				Required:    true,
				Description: DescModulePulumiFlagTask,
				Validators: []validator.String{
					stringvalidator.OneOf(pulumiCommandTaskValues...),
				},
			},
			"flag": schema.StringAttribute{
				Required:    true,
				Description: DescModulePulumiFlagFlag,
				Validators: []validator.String{
					stringvalidator.OneOf(pulumiFlagValues...),
				},
			},
			"value": schema.StringAttribute{
				Optional:    true,
				Description: DescModulePulumiFlagValue,
			},
		},
	}
}

func (r *modulePulumiFlagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data modulePulumiFlagModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(modulePulumiFlagEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *modulePulumiFlagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data modulePulumiFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", modulePulumiFlagEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *modulePulumiFlagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data modulePulumiFlagModel
	var state modulePulumiFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", modulePulumiFlagEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *modulePulumiFlagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data modulePulumiFlagModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", modulePulumiFlagEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *modulePulumiFlagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data modulePulumiFlagModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", modulePulumiFlagEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(modulePulumiFlagDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
