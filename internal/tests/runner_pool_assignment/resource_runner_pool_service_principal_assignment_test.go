// SPDX-License-Identifier: MPL-2.0

package runner_pool_assignment

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRunnerPoolServicePrincipalAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig +  core.RunnerPoolCreateConfig + RunnerPoolServicePrincipalAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_service_principal_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerPoolServicePrincipalAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig +  core.RunnerPoolCreateConfig + RunnerPoolServicePrincipalAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_service_principal_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_runner_pool_service_principal_assignment.this", "permission", "ReadWrite"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_runner_pool_service_principal_assignment" "this" { 
  runner_pool_id   = snapcd_runner_pool_service_principal_assignment.this.id
  service_principal_id   = data.snapcd_service_principal.this.id
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_service_principal_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_runner_pool_service_principal_assignment.this", "permission", "Read"),
				),
			},
		},
	})
}

func TestAccResourceRunnerPoolServicePrincipalAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig +  core.RunnerPoolCreateConfig + RunnerPoolServicePrincipalAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_service_principal_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_runner_pool_service_principal_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
