// SPDX-License-Identifier: MPL-2.0

package integration_events

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceNamespaceIntegrationEvent_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + IntegrationDataSourceConfig + NamespaceIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_integration_event.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_integration_event.this", "trigger", "JobFailed"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceIntegrationEvent_Update(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + IntegrationDataSourceConfig + NamespaceIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_integration_event.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + IntegrationDataSourceConfig + NamespaceIntegrationEventUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_integration_event.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_integration_event.this", "template", providerconfig.AppendRandomString("Job failed on namespace (updated): {{jobName}} %s")),
				),
			},
		},
	})
}

func TestAccResourceNamespaceIntegrationEvent_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig() + testdata.NamespaceCreateConfig + IntegrationDataSourceConfig + NamespaceIntegrationEventCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_integration_event.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_integration_event.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
