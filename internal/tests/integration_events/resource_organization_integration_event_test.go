// SPDX-License-Identifier: MPL-2.0

package integration_events

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceOrganizationIntegrationEvent_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + IntegrationDataSourceConfig + OrganizationIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_integration_event.this", "id"),
					resource.TestCheckResourceAttr("snapcd_organization_integration_event.this", "trigger", "MissionFaulted"),
				),
			},
		},
	})
}

func TestAccResourceOrganizationIntegrationEvent_Update(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + IntegrationDataSourceConfig + OrganizationIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_integration_event.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + IntegrationDataSourceConfig + OrganizationIntegrationEventUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_integration_event.this", "id"),
					resource.TestCheckResourceAttr("snapcd_organization_integration_event.this", "template", providerconfig.AppendRandomString("Mission faulted (updated): {{jobName}} %s")),
				),
			},
		},
	})
}

func TestAccResourceOrganizationIntegrationEvent_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + IntegrationDataSourceConfig + OrganizationIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_organization_integration_event.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_organization_integration_event.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
