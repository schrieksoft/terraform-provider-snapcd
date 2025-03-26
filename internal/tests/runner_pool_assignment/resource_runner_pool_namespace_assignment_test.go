// SPDX-License-Identifier: MPL-2.0

package runner_pool_assignment

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRunnerPoolNamespaceAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + core.RunnerPoolCreateConfig + RunnerPoolNamespaceAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_namespace_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerPoolNamespaceAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + core.RunnerPoolCreateConfig + RunnerPoolNamespaceAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_namespace_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_runner_pool_namespace_assignment.this", "permission", "ReadWrite"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_runner_pool_namespace_assignment" "this" { 
  runner_pool_id   = snapcd_runner_pool_namespace_assignment.this.id
  namespace_id   = snapcd_namespace.this.id
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_namespace_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_runner_pool_namespace_assignment.this", "permission", "Read"),
				),
			},
		},
	})
}

func TestAccResourceRunnerPoolNamespaceAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + core.RunnerPoolCreateConfig + RunnerPoolNamespaceAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_namespace_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_runner_pool_namespace_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
