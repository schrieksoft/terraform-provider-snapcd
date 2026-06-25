// SPDX-License-Identifier: MPL-2.0

package integration_events

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceModuleIntegrationEvent_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + IntegrationDataSourceConfig + ModuleIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_integration_event.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_integration_event.this", "trigger", "JobFailed"),
				),
			},
		},
	})
}

func TestAccResourceModuleIntegrationEvent_Update(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + IntegrationDataSourceConfig + ModuleIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_integration_event.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + IntegrationDataSourceConfig + ModuleIntegrationEventUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_integration_event.this", "id"),
					resource.TestCheckResourceAttr("snapcd_module_integration_event.this", "template", providerconfig.AppendRandomString("Job failed on module (updated): {{jobName}} %s")),
				),
			},
		},
	})
}

func TestAccResourceModuleIntegrationEvent_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + testdata.ModuleCreateConfig + IntegrationDataSourceConfig + ModuleIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_module_integration_event.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_module_integration_event.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
