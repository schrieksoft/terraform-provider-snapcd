// SPDX-License-Identifier: MPL-2.0

package integration_events

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStackIntegrationEvent_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.StackCreateConfig + IntegrationDataSourceConfig + StackIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_integration_event.this", "id"),
					resource.TestCheckResourceAttr("snapcd_stack_integration_event.this", "trigger", "JobFailed"),
				),
			},
		},
	})
}

func TestAccResourceStackIntegrationEvent_Update(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.StackCreateConfig + IntegrationDataSourceConfig + StackIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_integration_event.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.StackCreateConfig + IntegrationDataSourceConfig + StackIntegrationEventUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_integration_event.this", "id"),
					resource.TestCheckResourceAttr("snapcd_stack_integration_event.this", "template", providerconfig.AppendRandomString("Job failed on stack (updated): {{jobName}} %s")),
				),
			},
		},
	})
}

func TestAccResourceStackIntegrationEvent_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.StackCreateConfig + IntegrationDataSourceConfig + StackIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_stack_integration_event.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_stack_integration_event.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
