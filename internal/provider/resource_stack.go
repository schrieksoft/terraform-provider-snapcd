package provider

// import (
// 	"fmt"

// 	"context"

// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/types"

// 	snapcd "terraform-provider-snapcd/client"
// 	utils "terraform-provider-snapcd/utils"

// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// )

// var stackDefaultError = fmt.Sprintf("snapcd_stack error")

// var stackEndpoint = "/api/Definition/Stack"

// var _ resource.Resource = (*stackResource)(nil)

// func StackResource() resource.Resource {
// 	return &stackResource{}
// }

// type stackResource struct {
// 	client *snapcd.Client
// }

// // Configure adds the provider configured client to the resource.
// func (r *stackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}

// 	client, ok := req.ProviderData.(*snapcd.Client)

// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Unexpected Resource Configure Type",
// 			fmt.Sprintf("Expected *snapcd.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
// 		)

// 		return
// 	}

// 	r.client = client
// }

// func (r *stackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = req.ProviderTypeName + "_stack"
// }

// type stackModel struct {
// 	Name types.String `tfsdk:"name"`
// 	Id   types.String `tfsdk:"id"`
// }

// func (r *stackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		Attributes: map[string]schema.Attribute{
// 			"id": schema.StringAttribute{
// 				Computed: true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"name": schema.StringAttribute{
// 				Required: true,
// 			},
// 		},
// 	}
// }

// func (r *stackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// 	var data stackModel

// 	// Read Terraform plan data into the model
// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	jsonMap, err := utils.PlanToJson(data, []string{"id"})
// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Failed to convert json to plan: "+err.Error())
// 	}

// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Failed to convert plan to json: "+err.Error())
// 		return
// 	}

// 	result, httpError := r.client.Post(stackEndpoint, jsonMap)
// 	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
// 		resp.Diagnostics.AddError(globalRoleAssignmentDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
// 		return
// 	}
// 	if httpError != nil {
// 		err = httpError.Error
// 	} else {
// 		err = nil
// 	}
// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Error calling POST, unexpected error: "+err.Error())
// 		return
// 	}

// 	err = utils.JsonToPlan(result, &data)

// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Failed to convert json to plan: "+err.Error())
// 		return
// 	}

// 	// Save data into Terraform state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }

// func (r *stackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	var data stackModel

// 	// Read Terraform prior state data into the model
// 	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Read API call logic
// 	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", stackEndpoint, data.Id.ValueString()))
// 	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
// 		// Resource was not found, so remove it from state
// 		resp.State.RemoveResource(ctx)
// 		return
// 	}
// 	var err error
// 	if httpError != nil {
// 		err = httpError.Error
// 	} else {
// 		err = nil
// 	}
// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Error calling GET, unexpected error: "+err.Error())
// 		return
// 	}

// 	err = utils.JsonToPlan(result, &data)

// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Failed to convert json to plan: "+err.Error())
// 		return
// 	}

// 	// Save updated data into Terraform state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }

// func (r *stackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	var data stackModel
// 	var state stackModel

// 	// Read Terraform prior state data into the model
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

// 	// Read Terraform plan data into the model
// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Update API call logic

// 	jsonMap, err := utils.PlanToJson(data)
// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Failed to convert json to plan: "+err.Error())
// 	}

// 	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", stackEndpoint, state.Id.ValueString()), jsonMap)
// 	if httpError != nil {
// 		err = httpError.Error
// 	} else {
// 		err = nil
// 	}
// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Error calling PUT, unexpected error: "+err.Error())
// 		return
// 	}

// 	err = utils.JsonToPlan(result, &data)
// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Failed to convert json to plan: "+err.Error())
// 		return
// 	}

// 	// Save updated data into Terraform state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }

// func (r *stackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// 	var data stackModel

// 	// Read Terraform prior state data into the model
// 	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Delete API call logic
// 	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", stackEndpoint, data.Id.ValueString()))
// 	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
// 		// Resource was not found, so remove it from state
// 		resp.State.RemoveResource(ctx)
// 		return
// 	}
// 	var err error
// 	if httpError != nil {
// 		err = httpError.Error
// 	} else {
// 		err = nil
// 	}
// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
// 		return
// 	}
// }

// func (r *stackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	var data stackModel

// 	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", stackEndpoint, req.ID))
// 	var err error
// 	if httpError != nil {
// 		err = httpError.Error
// 	} else {
// 		err = nil
// 	}
// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Error calling GET, unexpected error: "+err.Error())
// 		return
// 	}

// 	err = utils.JsonToPlan(result, &data)
// 	if err != nil {
// 		resp.Diagnostics.AddError(stackDefaultError, "Failed to convert json to plan: "+err.Error())
// 		return
// 	}

// 	// Save updated data into Terraform state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }
