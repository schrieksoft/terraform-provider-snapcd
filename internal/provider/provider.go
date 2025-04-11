// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"strconv"

	snapcd "terraform-provider-snapcd/client"
	"terraform-provider-snapcd/internal/provider/core"
	"terraform-provider-snapcd/internal/provider/identity"
	"terraform-provider-snapcd/internal/provider/module_env_var"
	"terraform-provider-snapcd/internal/provider/module_param"
	"terraform-provider-snapcd/internal/provider/namespace_env_var"
	"terraform-provider-snapcd/internal/provider/namespace_param"
	"terraform-provider-snapcd/internal/provider/role_assignment"
	"terraform-provider-snapcd/internal/provider/secret"
	"terraform-provider-snapcd/internal/provider/secret_store"
	"terraform-provider-snapcd/internal/provider/secret_store_assignment"
	"terraform-provider-snapcd/internal/provider/runner_pool_assignment"
	"terraform-provider-snapcd/internal/provider/extra_files"
	"terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &snapcdProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &snapcdProvider{
			version: version,
		}
	}
}

// snapcdProvider is the provider implementation.
type snapcdProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// snapcdProviderModel maps provider schema data to a Go type.
type snapcdProviderModel struct {
	Url                        types.String `tfsdk:"url"`
	InsecureSkipVerify         types.Bool   `tfsdk:"insecure_skip_verify"`
	HealthCheckIntervalSeconds types.Int64  `tfsdk:"health_check_interval_seconds"`
	HealthCheckTimeoutSeconds  types.Int64  `tfsdk:"health_check_timeout_seconds"`
	AccessToken                types.String `tfsdk:"access_token"`
	ClientId                   types.String `tfsdk:"client_id"`
	ClientSecret               types.String `tfsdk:"client_secret"`
}

// Metadata returns the provider type name.
func (p *snapcdProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "snapcd"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *snapcdProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with snapcd.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "URL where the Snapcd API is served",
				Required:    true,
			},
			"access_token": schema.StringAttribute{
				Description: "Access token for the SnapCd API",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "Access token for the SnapCd API",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Access token for the SnapCd API",
				Optional:    true,
			},
			"insecure_skip_verify": schema.BoolAttribute{
				Description: "Skip verifying the HTTPs certificate of the API",
				Optional:    true,
			},
			"health_check_interval_seconds": schema.Int64Attribute{
				Description: "Number of seconds to wait inbetween polling the /health endpoint of the API. Defaults to 15 seconds.",
				Optional:    true,
			},
			"health_check_timeout_seconds": schema.Int64Attribute{
				Description: "Number of seconds during which to continuously poll /health endpoint before timing out. Defaults to 900 seconds (15 minutes)",
				Optional:    true,
			},
		},
	}
}

func (p *snapcdProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Snapcd client")

	// Retrieve provider data from configuration
	var config snapcdProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//////
	// insecure_skip_verify
	//////

	if config.InsecureSkipVerify.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("insecure_skip_verify"),
			"Unknown Snapcd InsecureSkipVerify",
			"The provider cannot create the Snapcd API client as there is an unknown configuration value for the Snapcd InsecureSkipVerify. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SNAPCD_INSECURE_SKIP_VERIFY environment variable.",
		)
	}

	insecure_skip_verify := false
	if !config.InsecureSkipVerify.IsNull() {
		insecure_skip_verify = config.InsecureSkipVerify.ValueBool()
	} else {
		insecure_skip_verify, _ = utils.GetEnvBool("SNAPCD_INSECURE_SKIP_VERIFY", false)
	}

	//////
	// url
	//////

	if config.Url.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown Snapcd Url",
			"The provider cannot create the Snapcd API client as there is an unknown configuration value for the Snapcd Url. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SNAPCD_URL environment variable.",
		)
	}

	url := os.Getenv("SNAPCD_URL")

	if !config.Url.IsNull() {
		url = config.Url.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing Snapcd API URL",
			"The provider cannot create the Snapcd API client as there is a missing or empty value for the Snapcd API URL. "+
				"Set the url value in the configuration or use the SNAPCD_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	//////
	// access_token
	//////

	if config.AccessToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_token"),
			"Unknown Snapcd AccessToken",
			"The provider cannot create the Snapcd API client as there is an unknown configuration value for the Snapcd AccessToken. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SNAPCD_ACCESS_TOKEN environment variable.",
		)
	}

	access_token := os.Getenv("SNAPCD_ACCESS_TOKEN")

	if !config.AccessToken.IsNull() {
		access_token = config.AccessToken.ValueString()
	}

	//////
	// client_id
	//////

	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown Snapcd AccessToken",
			"The provider cannot create the Snapcd API client as there is an unknown configuration value for the Snapcd AccessToken. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SNAPCD_CLIENT_ID environment variable.",
		)
	}

	client_id := os.Getenv("SNAPCD_CLIENT_ID")

	if !config.ClientId.IsNull() {
		client_id = config.ClientId.ValueString()
	}

	//////
	// client_secret
	//////

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret "),
			"Unknown Snapcd AccessToken",
			"The provider cannot create the Snapcd API client as there is an unknown configuration value for the Snapcd ClientSecret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SNAPCD_CLIENT_SECRET  environment variable.",
		)
	}

	client_secret := os.Getenv("SNAPCD_CLIENT_SECRET")

	if !config.ClientSecret.IsNull() {
		client_secret = config.ClientSecret.ValueString()
	}

	//////
	// null checks
	//////

	if access_token == "" && (client_id == "" || client_secret == "") {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_token"),
			"Missing Snapcd API authentication credentials",
			"The provider cannot create the Snapcd API client as there is a missing or empty value for both access_token and client_id / client_secret."+
				"Either of these must be set statically in the configuration, or using SNAPCD_ACCESS_TOKEN, or the combination of SNAPCD_CLIENT_ID/SNAPCD_CLIENT_SECRET. If either is already set, ensure the value is not empty.",
		)
	}

	//////
	// health_check_interval_seconds
	//////

	if config.HealthCheckIntervalSeconds.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("health_check_interval_seconds"),
			"Unknown Snapcd HealthCheckIntervalSeconds",
			"The provider cannot create the Snapcd API client as there is an unknown configuration value for the Snapcd HealthCheckIntervalSeconds. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SNAPCD_HEALTH_CHECK_INTERVAL_SECONDS environment variable.",
		)
	}

	health_check_interval_seconds_string, exists := os.LookupEnv("SNAPCD_HEALTH_CHECK_INTERVAL_SECONDS")

	if !exists {
		health_check_interval_seconds_string = "15"
	}

	// Try to convert it to an integer
	health_check_interval_seconds, err := strconv.Atoi(health_check_interval_seconds_string)

	// If there's an error converting the variable, print an error message and exit
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("health_check_interval_seconds"),
			"Ivalid Snapcd API HealthCheckIntervalSeconds",
			"The provider cannot create the Snapcd API client as there an invalid value for Snapcd API HealthCheckIntervalSeconds. "+
				"Set the health_check_interval_seconds value to a valid integer in the configuration or use the SNAPCD_HEALTH_CHECK_INTERVAL_SECONDS environment variable. "+
				"If either is already set, ensure the value is a valid integer.",
		)
	}

	if !config.HealthCheckIntervalSeconds.IsNull() {
		health_check_interval_seconds = int(config.HealthCheckIntervalSeconds.ValueInt64())
	}

	//////
	// health_check_timeout_seconds
	//////

	if config.HealthCheckTimeoutSeconds.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("health_check_timeout_seconds"),
			"Unknown Snapcd HealthCheckTimeoutSeconds",
			"The provider cannot create the Snapcd API client as there is an unknown configuration value for the Snapcd HealthCheckTimeoutSeconds. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SNAPCD_HEALTH_CHECK_TIMEOUT_SECONDS environment variable.",
		)
	}

	health_check_timeout_seconds_string, exists := os.LookupEnv("SNAPCD_HEALTH_CHECK_TIMEOUT_SECONDS")

	if !exists {
		health_check_timeout_seconds_string = "900"
	}

	// Try to convert it to an integer
	health_check_timeout_seconds, err := strconv.Atoi(health_check_timeout_seconds_string)

	// If there's an error converting the variable, print an error message and exit
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("health_check_timeout_seconds"),
			"Ivalid Snapcd API HealthCheckTimeoutSeconds",
			"The provider cannot create the Snapcd API client as there an invalid value for Snapcd API HealthCheckTimeoutSeconds. "+
				"Set the health_check_timeout_seconds value to a valid integer in the configuration or use the SNAPCD_HEALTH_CHECK_TIMEOUT_SECONDS environment variable. "+
				"If either is already set, ensure the value is a valid integer.",
		)
	}

	if !config.HealthCheckIntervalSeconds.IsNull() {
		health_check_timeout_seconds = int(config.HealthCheckTimeoutSeconds.ValueInt64())
	}

	// return if has error

	if resp.Diagnostics.HasError() {
		return
	}

	// set fields in context
	ctx = tflog.SetField(ctx, "insecure_skip_verify", insecure_skip_verify)
	ctx = tflog.SetField(ctx, "snapcd_url", url)
	ctx = tflog.SetField(ctx, "health_check_interval_seconds", health_check_interval_seconds)
	ctx = tflog.SetField(ctx, "health_check_timeout_seconds", health_check_timeout_seconds)
	ctx = tflog.SetField(ctx, "access_token", access_token)
	ctx = tflog.SetField(ctx, "client_id", client_id)
	ctx = tflog.SetField(ctx, "client_secret", client_secret)

	tflog.Debug(ctx, "Creating Snapcd client")

	// Create a new Snapcd client using the configuration values
	client, err := snapcd.CreateClient(url, insecure_skip_verify, health_check_interval_seconds, health_check_timeout_seconds, access_token, client_id, client_secret)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Snapcd API Client",
			"An unexpected error occurred when creating the Snapcd API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Snapcd Client Error: "+err.Error(),
		)
		return
	}

	// Make the Snapcd client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Snapcd client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *snapcdProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{

		core.NamespaceDataSource,
		core.ModuleDataSource,
		core.StackDataSource,
		core.RunnerPoolDataSource,

		identity.ServicePrincipalDataSource,
		identity.GroupDataSource,

		module_env_var.ModuleEnvVarFromLiteralDataSource,
		module_env_var.ModuleEnvVarFromDefinitionDataSource,
		module_env_var.ModuleEnvVarFromNamespaceDataSource,
		module_env_var.ModuleEnvVarFromOutputDataSource,
		module_env_var.ModuleEnvVarFromOutputSetDataSource,
		module_env_var.ModuleEnvVarFromSecretDataSource,

		module_param.ModuleParamFromLiteralDataSource,
		module_param.ModuleParamFromDefinitionDataSource,
		module_param.ModuleParamFromNamespaceDataSource,
		module_param.ModuleParamFromOutputDataSource,
		module_param.ModuleParamFromOutputSetDataSource,
		module_param.ModuleParamFromSecretDataSource,

		namespace_env_var.NamespaceEnvVarFromLiteralDataSource,
		namespace_env_var.NamespaceEnvVarFromDefinitionDataSource,
		namespace_env_var.NamespaceEnvVarFromSecretDataSource,

		namespace_param.NamespaceParamFromLiteralDataSource,
		namespace_param.NamespaceParamFromDefinitionDataSource,
		namespace_param.NamespaceParamFromSecretDataSource,

		secret_store.AzureKeyVaultSecretStoreDataSource,
		secret_store.SimpleSecretStoreDataSource,

		secret.AzureKeyVaultSecretScopedToStackDataSource,
		secret.AzureKeyVaultSecretScopedToNamespaceDataSource,
		secret.AzureKeyVaultSecretScopedToModuleDataSource,
		secret.AzureKeyVaultSecretDataSource,

		secret.SimpleSecretScopedToStackDataSource,
		secret.SimpleSecretScopedToNamespaceDataSource,
		secret.SimpleSecretScopedToModuleDataSource,
		secret.SimpleSecretDataSource,

		extra_files.NamespaceExtraFileDataSource,
		extra_files.ModuleExtraFileDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *snapcdProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{

		core.NamespaceResource,
		core.ModuleResource,
		core.StackResource,
		core.RunnerPoolResource,

		identity.ServicePrincipalResource,
		identity.GroupResource,
		identity.GroupMemberResource,

		module_env_var.ModuleEnvVarFromLiteralResource,
		module_env_var.ModuleEnvVarFromDefinitionResource,
		module_env_var.ModuleEnvVarFromNamespaceResource,
		module_env_var.ModuleEnvVarFromOutputResource,
		module_env_var.ModuleEnvVarFromOutputSetResource,
		module_env_var.ModuleEnvVarFromSecretResource,

		module_param.ModuleParamFromLiteralResource,
		module_param.ModuleParamFromDefinitionResource,
		module_param.ModuleParamFromNamespaceResource,
		module_param.ModuleParamFromOutputResource,
		module_param.ModuleParamFromOutputSetResource,
		module_param.ModuleParamFromSecretResource,

		namespace_env_var.NamespaceEnvVarFromLiteralResource,
		namespace_env_var.NamespaceEnvVarFromDefinitionResource,
		namespace_env_var.NamespaceEnvVarFromSecretResource,

		namespace_param.NamespaceParamFromLiteralResource,
		namespace_param.NamespaceParamFromDefinitionResource,
		namespace_param.NamespaceParamFromSecretResource,

		role_assignment.ResourceRoleAssignmentResource,
		role_assignment.TypeRoleAssignmentResource,
		role_assignment.GlobalRoleAssignmentResource,

		secret_store.AzureKeyVaultSecretStoreResource,
		secret_store.SimpleSecretStoreResource,

		secret.AzureKeyVaultSecretScopedToStackResource,
		secret.AzureKeyVaultSecretScopedToNamespaceResource,
		secret.AzureKeyVaultSecretScopedToModuleResource,
		secret.AzureKeyVaultSecretResource,

		secret.SimpleSecretScopedToStackResource,
		secret.SimpleSecretScopedToNamespaceResource,
		secret.SimpleSecretScopedToModuleResource,
		secret.SimpleSecretResource,

		secret_store_assignment.SecretStoreStackAssignmentResource,
		secret_store_assignment.SecretStoreNamespaceAssignmentResource,
		secret_store_assignment.SecretStoreModuleAssignmentResource,

		runner_pool_assignment.RunnerPoolStackAssignmentResource,
		runner_pool_assignment.RunnerPoolNamespaceAssignmentResource,
		runner_pool_assignment.RunnerPoolModuleAssignmentResource,

		extra_files.NamespaceExtraFileResource,
		extra_files.ModuleExtraFileResource,
	}
}
