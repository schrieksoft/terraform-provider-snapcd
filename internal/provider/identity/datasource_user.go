package identity

import (
	"fmt"

	"context"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = (*userDataSource)(nil)

var userDefaultError = fmt.Sprintf("snapcd_user error")

var userEndpoint = "/User"

func UserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *snapcd.Client
}

func (r *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

const (
	DescUserId               = "Unique ID of the User."
	DescIsDisabled           = "Indicates whether the user or not the user has been disabled."
	DescUserName             = "Unique name of the user."
	DescNormalizedUserName   = "Normalized user name used for consistency."
	DescEmail                = "User's email address."
	DescNormalizedEmail      = "Normalized email address used for consistency."
	DescEmailConfirmed       = "Whether the user's email has been confirmed."
	DescPasswordHash         = "Hashed password of the user."
	DescSecurityStamp        = "Security stamp used to identify changes to the user's security info."
	DescConcurrencyStamp     = "Used to handle concurrency checks."
	DescPhoneNumber          = "Phone number of the user."
	DescPhoneNumberConfirmed = "Whether the phone number is confirmed."
	DescTwoFactorEnabled     = "Indicates if two-factor authentication is enabled."
	DescLockoutEnd           = "The date and time when the lockout ends (if any)."
	DescLockoutEnabled       = "Indicates whether lockout is enabled for the user."
	DescAccessFailedCount    = "The number of failed access attempts."
)

type userModel struct {
	Id                   types.String `tfsdk:"id"`
	IsDisabled           types.Bool   `tfsdk:"is_disabled"`
	UserName             types.String `tfsdk:"user_name"`
	NormalizedUserName   types.String `tfsdk:"normalized_user_name"`
	Email                types.String `tfsdk:"email"`
	NormalizedEmail      types.String `tfsdk:"normalized_email"`
	EmailConfirmed       types.Bool   `tfsdk:"email_confirmed"`
	PasswordHash         types.String `tfsdk:"password_hash"`
	SecurityStamp        types.String `tfsdk:"security_stamp"`
	ConcurrencyStamp     types.String `tfsdk:"concurrency_stamp"`
	PhoneNumber          types.String `tfsdk:"phone_number"`
	PhoneNumberConfirmed types.Bool   `tfsdk:"phone_number_confirmed"`
	TwoFactorEnabled     types.Bool   `tfsdk:"two_factor_enabled"`
	LockoutEnd           types.String `tfsdk:"lockout_end"`
	LockoutEnabled       types.Bool   `tfsdk:"lockout_enabled"`
	AccessFailedCount    types.Int64  `tfsdk:"access_failed_count"`
}

func (d *userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Identity Access Management --- Use this data source to access information about an existing User in Snap CD.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: DescUserId,
			},
			"is_disabled": schema.BoolAttribute{
				Computed:    true,
				Description: DescIsDisabled,
			},
			"user_name": schema.StringAttribute{
				Required:    true,
				Description: DescUserName,
			},
			"normalized_user_name": schema.StringAttribute{
				Computed:    true,
				Description: DescNormalizedUserName,
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: DescEmail,
			},
			"normalized_email": schema.StringAttribute{
				Computed:    true,
				Description: DescNormalizedEmail,
			},
			"email_confirmed": schema.BoolAttribute{
				Computed:    true,
				Description: DescEmailConfirmed,
			},
			"password_hash": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: DescPasswordHash,
			},
			"security_stamp": schema.StringAttribute{
				Computed:    true,
				Description: DescSecurityStamp,
			},
			"concurrency_stamp": schema.StringAttribute{
				Computed:    true,
				Description: DescConcurrencyStamp,
			},
			"phone_number": schema.StringAttribute{
				Computed:    true,
				Description: DescPhoneNumber,
			},
			"phone_number_confirmed": schema.BoolAttribute{
				Computed:    true,
				Description: DescPhoneNumberConfirmed,
			},
			"two_factor_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: DescTwoFactorEnabled,
			},
			"lockout_end": schema.StringAttribute{
				Computed:    true,
				Description: DescLockoutEnd,
			},
			"lockout_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: DescLockoutEnabled,
			},
			"access_failed_count": schema.Int64Attribute{
				Computed:    true,
				Description: DescAccessFailedCount,
			},
		},
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, httpError := d.client.Get(fmt.Sprintf("%s/ByUserName/%s", userEndpoint, data.UserName.ValueString()))
	var err error
	if httpError != nil {
		err = httpError.Error
	} else {
		err = nil
	}

	if err != nil {
		resp.Diagnostics.AddError(userDefaultError, "Error creating calling GET, unexpected error: "+err.Error())
		return
	}

	err = utils.JsonToPlan(result, &data)

	if err != nil {
		resp.Diagnostics.AddError(userDefaultError, "Failed to convert map to struct: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
