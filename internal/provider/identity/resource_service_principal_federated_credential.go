package identity

import (
	"terraform-provider-snapcd/internal/provider/openapidocs"
	"terraform-provider-snapcd/internal/validators"

	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var servicePrincipalFederatedCredentialDefaultError = "snapcd_service_principal_federated_credential error"

var servicePrincipalFederatedCredentialEndpoint = "/ServicePrincipalFederatedCredential"

var _ resource.Resource = (*servicePrincipalFederatedCredentialResource)(nil)

func ServicePrincipalFederatedCredentialResource() resource.Resource {
	return &servicePrincipalFederatedCredentialResource{}
}

type servicePrincipalFederatedCredentialResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *servicePrincipalFederatedCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *servicePrincipalFederatedCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_principal_federated_credential"
}

type servicePrincipalFederatedCredentialModel struct {
	Id                 types.String `tfsdk:"id"`
	ServicePrincipalId types.String `tfsdk:"service_principal_id"`
	Name               types.String `tfsdk:"name"`
	Issuer             types.String `tfsdk:"issuer"`
	Subject            types.String `tfsdk:"subject"`
	Audience           types.String `tfsdk:"audience"`
	IsDisabled         types.Bool   `tfsdk:"is_disabled"`
}

func (r *servicePrincipalFederatedCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Identity Access Management --- Manages a Federated Credential on a Service Principal in Snap CD.

A workload presenting a token from the given issuer, carrying the given subject and audience, authenticates as the Service Principal. No client secret is involved, so nothing secret appears in this configuration or in Terraform state.

The Server must be able to reach the issuer to fetch its signing keys. For a Kubernetes ServiceAccount the subject is ` + "`system:serviceaccount:{namespace}:{name}`" + `.` + "\n\n## Required permissions\n\n" + openapidocs.ResourcePermissions["ServicePrincipalFederatedCredential"],
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: openapidocs.ServicePrincipalFederatedCredentialReadDto_Id,
			},
			"service_principal_id": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ServicePrincipalFederatedCredentialCreateDto_ServicePrincipalId,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ServicePrincipalFederatedCredentialCreateDto_Name,
				Validators: []validator.String{
					validators.NotBlank(),
				},
			},
			"issuer": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ServicePrincipalFederatedCredentialCreateDto_Issuer,
				Validators: []validator.String{
					validators.NotBlank(),
				},
			},
			"subject": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ServicePrincipalFederatedCredentialCreateDto_Subject,
				Validators: []validator.String{
					validators.NotBlank(),
				},
			},
			"audience": schema.StringAttribute{
				Required:    true,
				Description: openapidocs.ServicePrincipalFederatedCredentialCreateDto_Audience,
				Validators: []validator.String{
					validators.NotBlank(),
				},
			},
			"is_disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: openapidocs.ServicePrincipalFederatedCredentialCreateDto_IsDisabled,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *servicePrincipalFederatedCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data servicePrincipalFederatedCredentialModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(servicePrincipalFederatedCredentialEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *servicePrincipalFederatedCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data servicePrincipalFederatedCredentialModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", servicePrincipalFederatedCredentialEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		utils.ExternallyDeleted(ctx, &resp.Diagnostics, "snapcd_service_principal_federated_credential", data.Id.ValueString())
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
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *servicePrincipalFederatedCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data servicePrincipalFederatedCredentialModel
	var state servicePrincipalFederatedCredentialModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data)
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", servicePrincipalFederatedCredentialEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *servicePrincipalFederatedCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data servicePrincipalFederatedCredentialModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", servicePrincipalFederatedCredentialEndpoint, data.Id.ValueString()))
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
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *servicePrincipalFederatedCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data servicePrincipalFederatedCredentialModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", servicePrincipalFederatedCredentialEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(servicePrincipalFederatedCredentialDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
