package provider

import (
	snapcd "terraform-provider-snapcd/client"

	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"		
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

)

// Concrete Resource Implementations
// Stack Resource
type stackModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (m *stackModel) GetID() types.String   { return m.Id }
func (m *stackModel) SetID(id types.String) { m.Id = id }

var stackEndpoint = "/api/Definition/Stack"

func StackResource() resource.Resource {
	return NewResource[*stackModel](ResourceConfig{
		TypeName:           "stack",
		Endpoint:           "/api/Definition/Stack",
		DefaultError:       "snapcd_stack error",
		EntityNotFoundCode: snapcd.Status441EntityNotFound,
		EntityExistsCode:   snapcd.Status442EntityAlreadyExists,
		Schema: resourceSchema.Schema{
			Attributes: map[string]resourceSchema.Attribute{
				"id": resourceSchema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"name": resourceSchema.StringAttribute{
					Required: true,
				},
			},
		},
	})
}


func StackDataSource() datasource.DataSource {
	return NewDataSource[*stackModel](DataSourceConfig{
		TypeName:     "stack",
		DefaultError: "snapcd_stack error",
		Schema: datasourceSchema.Schema{
			Attributes: map[string]datasourceSchema.Attribute{
				"id": datasourceSchema.StringAttribute{
					Computed: true,
				},
				"name": datasourceSchema.StringAttribute{
					Required: true,
				},
			},
		},
		EndpointFunc: func(model any) string {
			m := model.(*stackModel)
			return fmt.Sprintf("%s/ByName/%s", stackEndpoint, m.Name.ValueString())
		},
	})
}
