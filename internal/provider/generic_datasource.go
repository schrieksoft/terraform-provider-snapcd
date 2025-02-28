package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"
)

// DataSource Configuration
type DataSourceConfig struct {
	TypeName     string
	DefaultError string
	Schema       schema.Schema
	EndpointFunc func(model any) string // Flexible endpoint calculation
}

// Model Interface
type DataSourceModel interface {
	GetID() types.String
	SetID(types.String)
}

// Generic DataSource Implementation
type genericDataSource[T DataSourceModel] struct {
	client *snapcd.Client
	config DataSourceConfig
}

func NewDataSource[T DataSourceModel](config DataSourceConfig) datasource.DataSource {
	return &genericDataSource[T]{config: config}
}

func (d *genericDataSource[T]) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.config.TypeName
}

func (d *genericDataSource[T]) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.config.Schema
}

func (d *genericDataSource[T]) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*snapcd.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *snapcd.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *genericDataSource[T]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model T

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := d.config.EndpointFunc(model)
	result, httpErr := d.client.Get(endpoint)
	if httpErr != nil {
		resp.Diagnostics.AddError(d.config.DefaultError, "API Error: "+httpErr.Error.Error())
		return
	}

	if err := utils.JsonToPlan(result, &model); err != nil {
		resp.Diagnostics.AddError(d.config.DefaultError, "Deserialization error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
