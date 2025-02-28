package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	snapcd "terraform-provider-snapcd/client"
	utils "terraform-provider-snapcd/utils"
)

// Resource Configuration
type ResourceConfig struct {
	TypeName           string
	Endpoint           string
	DefaultError       string
	EntityNotFoundCode int
	EntityExistsCode   int
	Schema             schema.Schema
}

// Model Interface
type ResourceModel interface {
	GetID() types.String
	SetID(types.String)
}

// Generic Resource Implementation
type genericResource[T ResourceModel] struct {
	client *snapcd.Client
	config ResourceConfig
}

func NewResource[T ResourceModel](config ResourceConfig) resource.Resource {
	return &genericResource[T]{config: config}
}

func (r *genericResource[T]) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.config.TypeName
}

func (r *genericResource[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.config.Schema
}

func (r *genericResource[T]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*snapcd.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *snapcd.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = client
}

// Create Method
func (r *genericResource[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan T
	plan = *new(T) // Initialize the concrete type

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(plan, []string{"id"})
	if err != nil {
		resp.Diagnostics.AddError(r.config.DefaultError, "Serialization error: "+err.Error())
		return
	}

	result, httpErr := r.client.Post(r.config.Endpoint, jsonMap)
	if httpErr != nil {
		if httpErr.StatusCode == r.config.EntityExistsCode {
			resp.Diagnostics.AddError(r.config.DefaultError, "Resource exists - must import")
			return
		}
		resp.Diagnostics.AddError(r.config.DefaultError, "API Error: "+httpErr.Error.Error())
		return
	}

	if err := utils.JsonToPlan(result, &plan); err != nil {
		resp.Diagnostics.AddError(r.config.DefaultError, "Deserialization error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read Method
func (r *genericResource[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state T
	state = *new(T) // Initialize the concrete type

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/%s", r.config.Endpoint, state.GetID().ValueString())
	result, httpErr := r.client.Get(url)
	if httpErr != nil {
		if httpErr.StatusCode == r.config.EntityNotFoundCode {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(r.config.DefaultError, "API Error: "+httpErr.Error.Error())
		return
	}

	if err := utils.JsonToPlan(result, &state); err != nil {
		resp.Diagnostics.AddError(r.config.DefaultError, "Deserialization error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update Method
func (r *genericResource[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state T
	plan = *new(T) // Initialize the concrete type
	state = *new(T)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jsonMap, err := utils.PlanToJson(plan, nil)
	if err != nil {
		resp.Diagnostics.AddError(r.config.DefaultError, "Serialization error: "+err.Error())
		return
	}

	url := fmt.Sprintf("%s/%s", r.config.Endpoint, state.GetID().ValueString())
	result, httpErr := r.client.Put(url, jsonMap)
	if httpErr != nil {
		resp.Diagnostics.AddError(r.config.DefaultError, "API Error: "+httpErr.Error.Error())
		return
	}

	if err := utils.JsonToPlan(result, &plan); err != nil {
		resp.Diagnostics.AddError(r.config.DefaultError, "Deserialization error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete Method
func (r *genericResource[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state T
	state = *new(T) // Initialize the concrete type

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/%s", r.config.Endpoint, state.GetID().ValueString())
	if _, httpErr := r.client.Delete(url); httpErr != nil {
		if httpErr.StatusCode != r.config.EntityNotFoundCode {
			resp.Diagnostics.AddError(r.config.DefaultError, "API Error: "+httpErr.Error.Error())
		}
	}
}

// ImportState Method
func (r *genericResource[T]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var model T
	model = *new(T) // Initialize the concrete type

	url := fmt.Sprintf("%s/%s", r.config.Endpoint, req.ID)
	result, httpErr := r.client.Get(url)
	if httpErr != nil {
		resp.Diagnostics.AddError(r.config.DefaultError, "Import Error: "+httpErr.Error.Error())
		return
	}

	if err := utils.JsonToPlan(result, &model); err != nil {
		resp.Diagnostics.AddError(r.config.DefaultError, "Import Deserialization Error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// Utility Functions
func PlanToJson(plan interface{}, excludeFields []string) (map[string]interface{}, error) {
	val := reflect.ValueOf(plan)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to struct")
	}

	result := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if !shouldExclude(field.Name, excludeFields) {
			result[field.Name] = val.Field(i).Interface()
		}
	}

	return result, nil
}

func shouldExclude(fieldName string, excludeFields []string) bool {
	for _, f := range excludeFields {
		if f == fieldName {
			return true
		}
	}
	return false
}
