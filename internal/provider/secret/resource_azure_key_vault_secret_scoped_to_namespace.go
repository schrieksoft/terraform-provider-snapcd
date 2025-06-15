package secret

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

var azureKeyVaultSecretScopedToNamespaceDefaultError = fmt.Sprintf("snapcd_azure_key_vault_secret_scoped_to_namespace error")

var azureKeyVaultSecretScopedToNamespaceEndpoint = "/api/AzureKeyVaultSecretScopedToNamespace"

var _ resource.Resource = (*azureKeyVaultSecretScopedToNamespaceResource)(nil)

func AzureKeyVaultSecretScopedToNamespaceResource() resource.Resource {
	return &azureKeyVaultSecretScopedToNamespaceResource{}
}

type azureKeyVaultSecretScopedToNamespaceResource struct {
	client *snapcd.Client
}

// Configure adds the provider configured client to the resource.
func (r *azureKeyVaultSecretScopedToNamespaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *azureKeyVaultSecretScopedToNamespaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_key_vault_secret_scoped_to_namespace"
}

type azureKeyVaultSecretScopedToNamespaceModel struct {
	Name             types.String `tfsdk:"name"`
	Id               types.String `tfsdk:"id"`
	SecretStoreId    types.String `tfsdk:"secret_store_id"`
	NamespaceId      types.String `tfsdk:"namespace_id"`
	RemoteSecretName types.String `tfsdk:"remote_secret_name"`
}

func (r *azureKeyVaultSecretScopedToNamespaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Secrets --- Manages a Azure Key Vault Secret (Scoped to Namespace) in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: DescId,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: DescName,
			},
			"remote_secret_name": schema.StringAttribute{
				Required:    true,
				Description: DescRemoteName,
			},
			"namespace_id": schema.StringAttribute{
				Required:    true,
				Description: DescNamespaceId,
			},
			"secret_store_id": schema.StringAttribute{
				Required:    true,
				Description: DescSecretStoreId,
			},
		},
	}
}

func (r *azureKeyVaultSecretScopedToNamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data azureKeyVaultSecretScopedToNamespaceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(data, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Failed to convert plan to json: "+err.Error())
		return
	}

	result, httpError := r.client.Post(azureKeyVaultSecretScopedToNamespaceEndpoint, jsonMap)
	if httpError != nil && httpError.StatusCode == snapcd.Status442EntityAlreadyExists {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "The resource you are trying to create already exists. To manage it with terraform you must import it")
		return
	}
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Error calling POST, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *azureKeyVaultSecretScopedToNamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data azureKeyVaultSecretScopedToNamespaceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", azureKeyVaultSecretScopedToNamespaceEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		// Resource was not found, so remove it from state
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
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *azureKeyVaultSecretScopedToNamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data azureKeyVaultSecretScopedToNamespaceModel
	var state azureKeyVaultSecretScopedToNamespaceModel

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
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Failed to convert json to plan: "+err.Error())
	}

	result, httpError := r.client.Put(fmt.Sprintf("%s/%s", azureKeyVaultSecretScopedToNamespaceEndpoint, state.Id.ValueString()), jsonMap)
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Error calling PUT, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *azureKeyVaultSecretScopedToNamespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data azureKeyVaultSecretScopedToNamespaceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, httpError := r.client.Delete(fmt.Sprintf("%s/%s", azureKeyVaultSecretScopedToNamespaceEndpoint, data.Id.ValueString()))
	if httpError != nil && httpError.StatusCode == snapcd.Status441EntityNotFound {
		// Resource was not found, so remove it from state
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
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Error calling DELETE, unexpected error: "+err.Error())
		return
	}
}

func (r *azureKeyVaultSecretScopedToNamespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data azureKeyVaultSecretScopedToNamespaceModel

	result, httpError := r.client.Get(fmt.Sprintf("%s/%s", azureKeyVaultSecretScopedToNamespaceEndpoint, req.ID))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Error calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)
	if err != nil {
		resp.Diagnostics.AddError(azureKeyVaultSecretScopedToNamespaceDefaultError, "Failed to convert json to plan: "+err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
